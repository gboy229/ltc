// Package menu 菜单注册与同步
// 菜单定义从 SQL 迁移到 Go 代码，通过 menu 命令同步到 hg_admin_menu 表。
// 只定义「目录」和「菜单」，按钮由操作类型(Ops)自动生成。
package menu

import (
	"fmt"
	"strings"
)

// 菜单类型
const (
	TypeDir  = 1 // 目录
	TypeMenu = 2 // 菜单
	TypeBtn  = 3 // 按钮
)

// Button 自定义按钮
type Button struct {
	Name        string // 按钮唯一标识
	Title       string // 按钮名称
	Permissions string // 权限路由
}

// Item 菜单定义（用结构体字面量直接表达菜单树）
type Item struct {
	Name       string    // 唯一标识（同步键），如 "org_dept"
	Title      string    // 菜单名称，如 "部门管理"
	Path       string    // 前端路由地址
	Icon       string    // 菜单图标（目录用）
	Component  string    // 组件路径（菜单用）
	Api        string    // API 前缀，用于生成读权限(/api/list)和按钮权限(/api/edit 等)
	Perm       string    // 读权限覆盖（可选，默认自动推导 /api/list）
	Ops        string    // 操作类型，逗号分隔："edit,delete,status"
	Buttons    []*Button // 自定义按钮（标准操作之外）
	Sort       int       // 排序
	Dir        bool      // 目录标记（默认菜单）
	Children   []*Item   // 子菜单
	Redirect   string    // 重定向地址
	ActiveMenu string    // 高亮菜单编码
	FrameSrc   string    // 内联外部地址
	KeepAlive  bool      // 缓存路由
	Hidden     bool      // 是否隐藏
	Affix      bool      // 是否固定
	AlwaysShow bool      // 取消自动计算根路由模式
	IsRoot     bool      // 是否根路由
	IsFrame    bool      // 是否内嵌
}

// opDef 操作类型 → 按钮约定
type opDef struct {
	Suffix   string // 按钮名后缀
	TitleTpl string // 标题模板
	Route    string // 权限路由后缀
}

// opMap 操作类型映射，对齐代码生成模板 source.sql.template
var opMap = map[string]opDef{
	"view":       {Suffix: "View", TitleTpl: "%s详情", Route: "/view"},
	"edit":       {Suffix: "Edit", TitleTpl: "编辑/新增%s", Route: "/edit"},
	"maxSort":    {Suffix: "MaxSort", TitleTpl: "获取%s最大排序", Route: "/maxSort"},
	"delete":     {Suffix: "Delete", TitleTpl: "删除%s", Route: "/delete"},
	"status":     {Suffix: "Status", TitleTpl: "修改%s状态", Route: "/status"},
	"switch":     {Suffix: "Switch", TitleTpl: "操作%s开关", Route: "/switch"},
	"export":     {Suffix: "Export", TitleTpl: "导出%s", Route: "/export"},
	"treeOption": {Suffix: "TreeOption", TitleTpl: "获取%s关系树选项", Route: "/treeOption"},
}

var (
	items       []*Item                    // 底座菜单
	pluginItems = make(map[string][]*Item) // 插件菜单（scope → items）
)

// Register 注册底座顶级菜单
func Register(list ...*Item) {
	items = append(items, list...)
}

// RegisterPlugin 注册插件菜单
func RegisterPlugin(scope string, list ...*Item) {
	pluginItems[scope] = append(pluginItems[scope], list...)
}

// Items 获取所有已注册的底座顶级菜单
func Items() []*Item {
	return items
}

// GenButtons 根据操作类型生成按钮
func (i *Item) GenButtons() []*Button {
	var btns []*Button
	for _, op := range strings.Split(i.Ops, ",") {
		op = strings.TrimSpace(op)
		if op == "" {
			continue
		}
		def, ok := opMap[op]
		if !ok {
			continue
		}
		btns = append(btns, &Button{
			Name:        i.apiPrefix() + "_" + op,
			Title:       fmt.Sprintf(def.TitleTpl, i.Title),
			Permissions: "/" + i.apiPrefix() + def.Route,
		})
	}
	btns = append(btns, i.Buttons...)
	return btns
}

// apiPrefix 返回 API 前缀：Api 优先，否则默认取 Path
func (i *Item) apiPrefix() string {
	if i.Api != "" {
		return i.Api
	}
	return i.Path
}

// readPermissions 读权限：Perm 优先；有操作(Api/Ops)时才自动推导 /api/list；否则为空
func (i *Item) readPermissions() string {
	if i.Perm != "" {
		return i.Perm
	}
	if i.Api == "" && i.Ops == "" {
		return ""
	}
	return "/" + i.apiPrefix() + "/list"
}
