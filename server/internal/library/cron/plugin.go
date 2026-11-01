// Package cron
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package cron

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
)

// CronItem 插件定时任务定义
type CronItem struct {
	Name    string // 任务方法名（对应 cron.Register 的 name，带插件前缀）
	Title   string // 任务标题
	Pattern string // 默认表达式
	Policy  int64  // 策略：1并行 2单例 3单次 4多次
	GroupId int64  // 分组ID
	Params  string // 函数参数
	Sort    int    // 排序
	Remark  string // 备注
	Status  int    // 状态
}

var pluginCronItems = make(map[string][]*CronItem)

// RegisterPlugin 注册插件定时任务定义
func RegisterPlugin(scope string, items ...*CronItem) {
	pluginCronItems[scope] = append(pluginCronItems[scope], items...)
}

// SyncPlugin 同步插件定时任务（安装/更新，按 name 前缀 upsert）
func SyncPlugin(ctx context.Context, scope string) (err error) {
	cols := dao.SysCron.Columns()

	// 现有插件任务（按 name 前缀匹配）
	var existing []*entity.SysCron
	if err = dao.SysCron.Ctx(ctx).Where(cols.Name+" LIKE ?", scope+"\\_%").Scan(&existing); err != nil {
		return
	}
	existingMap := make(map[string]*entity.SysCron, len(existing))
	for _, c := range existing {
		existingMap[c.Name] = c
	}

	synced := make(map[string]bool)
	for _, it := range pluginCronItems[scope] {
		data := g.Map{
			cols.Name:    it.Name,
			cols.Title:   it.Title,
			cols.Pattern: it.Pattern,
			cols.Policy:  it.Policy,
			cols.GroupId: it.GroupId,
			cols.Params:  it.Params,
			cols.Sort:    it.Sort,
			cols.Remark:  it.Remark,
			cols.Status:  it.Status,
		}
		if em, ok := existingMap[it.Name]; ok {
			if _, err = dao.SysCron.Ctx(ctx).WherePri(em.Id).Data(data).Update(); err != nil {
				return
			}
		} else {
			if _, err = dao.SysCron.Ctx(ctx).Data(data).Insert(); err != nil {
				return
			}
		}
		synced[it.Name] = true
	}

	// 删除插件内不再声明的任务
	for name, c := range existingMap {
		if !synced[name] {
			if _, err = dao.SysCron.Ctx(ctx).WherePri(c.Id).Delete(); err != nil {
				return
			}
		}
	}
	return
}

// DeletePlugin 删除插件定时任务（卸载）
func DeletePlugin(ctx context.Context, scope string) (err error) {
	_, err = dao.SysCron.Ctx(ctx).Where(dao.SysCron.Columns().Name+" LIKE ?", scope+"\\_%").Delete()
	return
}
