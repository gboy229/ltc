// Package menu 菜单同步
package menu

import (
	"context"
	"strings"

	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	treeutil "hotgo/utility/tree"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// statusDisable 与 internal/statusDisable 保持一致；此处不 import consts，避免循环依赖
const statusDisable = 2

// Sync 将注册的底座菜单同步到数据库
func Sync(ctx context.Context) (err error) {
	return syncScope(ctx, "", Items())
}

// SyncPlugin 将指定插件的菜单同步到数据库
// 以菜单 name 作为同步键：存在则更新，不存在则新增，ID 由数据库自增分配
func SyncPlugin(ctx context.Context, scope string) (err error) {
	return syncScope(ctx, scope, pluginItems[scope])
}

// DisablePlugin 禁用指定插件的所有菜单（卸载时调用）
// 不删除菜单行，重装时按 name 复用原 ID，避免角色-菜单权限关系失效
func DisablePlugin(ctx context.Context, scope string) (err error) {
	cols := dao.AdminMenu.Columns()
	result, err := dao.AdminMenu.Ctx(ctx).
		Where(cols.Scope, scope).
		Data(g.Map{cols.Status: statusDisable}).
		Update()
	if err != nil {
		return gerror.Wrapf(err, "禁用插件菜单失败: %s", scope)
	}
	n, _ := result.RowsAffected()
	g.Log().Infof(ctx, "禁用插件菜单: %s，共 %d 个", scope, n)
	return nil
}

// DeletePlugin Deprecated: 卸载请用 DisablePlugin。
// 为兼容旧插件和脚手架模板，这里不再硬删除，改为软禁用，避免菜单 ID 变化导致角色权限失效。
func DeletePlugin(ctx context.Context, scope string) (err error) {
	return DisablePlugin(ctx, scope)
}

// PurgePlugin 永久删除指定插件的菜单及角色关联（仅用于显式清理）
func PurgePlugin(ctx context.Context, scope string) (err error) {
	cols := dao.AdminMenu.Columns()
	ids, err := dao.AdminMenu.Ctx(ctx).Where(cols.Scope, scope).Fields(cols.Id).Array()
	if err != nil {
		return gerror.Wrapf(err, "查询插件菜单失败: %s", scope)
	}
	if len(ids) > 0 {
		roleMenuCols := dao.AdminRoleMenu.Columns()
		if _, err = dao.AdminRoleMenu.Ctx(ctx).WhereIn(roleMenuCols.MenuId, ids).Delete(); err != nil {
			return gerror.Wrapf(err, "删除角色菜单关联失败: %s", scope)
		}
	}
	result, err := dao.AdminMenu.Ctx(ctx).Where(cols.Scope, scope).Delete()
	if err != nil {
		return gerror.Wrapf(err, "删除插件菜单失败: %s", scope)
	}
	n, _ := result.RowsAffected()
	g.Log().Infof(ctx, "永久删除插件菜单: %s，共 %d 个", scope, n)
	return nil
}

// syncScope 按作用域同步菜单（scope 为空表示底座），顶级项 pid=0
func syncScope(ctx context.Context, scope string, items []*Item) (err error) {
	cols := dao.AdminMenu.Columns()

	// 加载现有菜单（按作用域过滤）
	var existing []*entity.AdminMenu
	if err = dao.AdminMenu.Ctx(ctx).Where(cols.Scope, scope).Scan(&existing); err != nil {
		return
	}
	existingMap := make(map[string]*entity.AdminMenu, len(existing))
	for _, m := range existing {
		existingMap[m.Name] = m
	}

	syncedNames := make(map[string]bool)

	// 递归同步（先父后子）
	var syncItem func(it *Item, pid int64, pLevel int, pTree string, parentPath string, sort int) error
	syncItem = func(it *Item, pid int64, pLevel int, pTree string, parentPath string, sort int) error {
		level := 1
		tree := ""
		if pid > 0 {
			level = pLevel + 1
			tree = treeutil.GenLabel(pTree, pid)
		}

		typ := TypeMenu
		if it.Dir {
			typ = TypeDir
		}

		// 计算完整路径
		fullPath := it.Path
		if parentPath != "" && !strings.HasPrefix(it.Path, "/") {
			fullPath = parentPath + "/" + it.Path
		}

		// 构造写入项（应用默认值，避免修改原注册项）
		w := *it
		if it.Sort > 0 {
			w.Sort = it.Sort // 显式设置过排序时优先使用
		} else {
			w.Sort = sort // 否则按注册顺序自动分配
		}
		if it.Dir {
			w.AlwaysShow = true // 目录默认始终显示
			if w.Redirect == "" && len(it.Children) > 0 {
				w.Redirect = fullPath + "/" + it.Children[0].Path // 默认重定向到第一个子菜单
			}
			if w.Component == "" {
				if pid == 0 {
					w.Component = "LAYOUT" // 顶级目录：主布局
				} else {
					w.Component = "ParentLayout" // 嵌套目录：父子布局
				}
			}
		}

		// 仅菜单(type=2)有读权限，目录/按钮无
		permissions := ""
		if typ == TypeMenu {
			permissions = it.readPermissions()
		}

		newId, e := upsert(ctx, &w, typ, permissions, pid, level, tree, scope, existingMap)
		if e != nil {
			return e
		}
		syncedNames[it.Name] = true

		// 菜单(type=2)自动生成按钮(type=3)，按钮是菜单的子级
		if typ == TypeMenu {
			btnLevel := level + 1
			btnTree := treeutil.GenLabel(tree, newId)
			for _, btn := range it.GenButtons() {
				btnItem := &Item{Name: btn.Name, Title: btn.Title, Sort: 10}
				if _, e = upsert(ctx, btnItem, TypeBtn, btn.Permissions, newId, btnLevel, btnTree, scope, existingMap); e != nil {
					return e
				}
				syncedNames[btn.Name] = true
			}
		}

		// 递归子菜单（排序按注册顺序）
		for i, child := range it.Children {
			if e = syncItem(child, newId, level, tree, fullPath, (i+1)*10); e != nil {
				return e
			}
		}
		return nil
	}

	for i, it := range items {
		if err = syncItem(it, 0, 0, "", "", (i+1)*10); err != nil {
			return
		}
	}

	// 禁用作用域内不在代码里的菜单，不删除行，保留 ID 和角色权限关系
	disabled := 0
	for name, m := range existingMap {
		if !syncedNames[name] {
			if _, err = dao.AdminMenu.Ctx(ctx).WherePri(m.Id).Data(g.Map{
				cols.Status: statusDisable,
			}).Update(); err != nil {
				return gerror.Wrapf(err, "禁用菜单失败: %s", name)
			}
			disabled++
			g.Log().Infof(ctx, "禁用菜单: %s", name)
		}
	}

	g.Log().Infof(ctx, "菜单同步完成（scope=%q）：共 %d 个节点，禁用 %d 个", scope, len(syncedNames), disabled)
	return nil
}

// upsert 按 name 新增或更新单个菜单节点，返回其 ID
func upsert(ctx context.Context, it *Item, typ int, permissions string, pid int64, level int, tree string, scope string, existingMap map[string]*entity.AdminMenu) (id int64, err error) {
	cols := dao.AdminMenu.Columns()

	data := g.Map{
		cols.Pid:         pid,
		cols.Level:       level,
		cols.Tree:        tree,
		cols.Name:        it.Name,
		cols.Title:       it.Title,
		cols.Type:        typ,
		cols.Scope:       scope,
		cols.Path:        it.Path,
		cols.Icon:        it.Icon,
		cols.Component:   it.Component,
		cols.Redirect:    it.Redirect,
		cols.Permissions: permissions,
		cols.Sort:        it.Sort,
		cols.Status:      1,
	}
	appendBool(data, cols.KeepAlive, it.KeepAlive)
	appendBool(data, cols.Hidden, it.Hidden)
	appendBool(data, cols.Affix, it.Affix)
	appendBool(data, cols.AlwaysShow, it.AlwaysShow)
	appendBool(data, cols.IsFrame, it.IsFrame)
	appendBool(data, cols.IsRoot, it.IsRoot)
	if it.ActiveMenu != "" {
		data[cols.ActiveMenu] = it.ActiveMenu
	}
	if it.FrameSrc != "" {
		data[cols.FrameSrc] = it.FrameSrc
	}

	if em, ok := existingMap[it.Name]; ok {
		// 修改（保留原 ID）
		id = em.Id
		if _, err = dao.AdminMenu.Ctx(ctx).WherePri(id).Data(data).Update(); err != nil {
			return 0, gerror.Wrapf(err, "更新菜单失败: %s", it.Name)
		}
		g.Log().Debugf(ctx, "更新菜单: %s", it.Name)
	} else {
		// 新增：ID 由数据库自增分配
		data[cols.CreatedAt] = gtime.Now()
		data[cols.UpdatedAt] = gtime.Now()
		if id, err = dao.AdminMenu.Ctx(ctx).Data(data).InsertAndGetId(); err != nil {
			return 0, gerror.Wrapf(err, "新增菜单失败: %s", it.Name)
		}
		g.Log().Infof(ctx, "新增菜单: %s", it.Name)
	}
	return id, nil
}

// appendBool 仅在值为 true 时写入 1，避免覆盖数据库默认值
func appendBool(data g.Map, field string, val bool) {
	if val {
		data[field] = 1
	}
}
