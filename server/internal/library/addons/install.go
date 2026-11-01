// Package addons
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package addons

import (
	"context"
	"io/fs"
	"path"

	addonsFS "hotgo/addons"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/migrate"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
)

// InstallRecord 安装记录
type InstallRecord struct {
	Id        int64       `json:"id"        description:"安装ID"`
	Version   string      `json:"version"   description:"安装版本"`
	Status    int         `json:"status"    description:"安装状态"`
	CreatedAt *gtime.Time `json:"createdAt" description:"创建时间"`
	UpdatedAt *gtime.Time `json:"updatedAt" description:"更新时间"`
}

func GetModel(ctx context.Context) *gdb.Model {
	return dao.SysAddonsInstall.Ctx(ctx)
}

func ScanInstall(m Module) (record *InstallRecord, err error) {
	err = GetModel(m.Ctx()).Where(dao.SysAddonsInstall.Columns().Name, m.GetSkeleton().Name).Scan(&record)
	return
}

// IsInstall 模块是否已安装
func IsInstall(m Module) bool {
	record, err := ScanInstall(m)
	if err != nil {
		g.Log().Debugf(m.Ctx(), "addons.IsInstall err:%+v", err)
		return false
	}
	if record == nil {
		return false
	}
	return record.Status == consts.AddonsInstallStatusOk
}

// Install 安装模块
func Install(m Module) (err error) {
	record, err := ScanInstall(m)
	if err != nil {
		return
	}
	if record != nil && record.Status == consts.AddonsInstallStatusOk {
		return gerror.New("插件已安装，无需重复操作！")
	}

	// 执行插件结构迁移（建表），幂等
	if err = migrateUp(m.Ctx(), m.GetSkeleton().Name); err != nil {
		return
	}

	// 执行插件数据填充（种子），幂等
	if err = migrateSeed(m.Ctx(), m.GetSkeleton().Name); err != nil {
		return
	}

	data := g.Map{
		dao.SysAddonsInstall.Columns().Name:    m.GetSkeleton().Name,
		dao.SysAddonsInstall.Columns().Version: m.GetSkeleton().Version,
		dao.SysAddonsInstall.Columns().Status:  consts.AddonsInstallStatusOk,
	}
	return g.DB().Transaction(m.Ctx(), func(ctx context.Context, tx gdb.TX) error {
		if record != nil {
			_, _ = GetModel(ctx).Where(dao.SysAddonsInstall.Columns().Id, record.Id).Delete()
		}

		if _, err = GetModel(ctx).Data(data).OmitEmptyData().Insert(); err != nil {
			return err
		}
		return m.Install(ctx)
	})
}

// Upgrade 更新模块
func Upgrade(m Module) (err error) {
	record, err := ScanInstall(m)
	if err != nil {
		return
	}

	if record == nil || record.Status != consts.AddonsInstallStatusOk {
		return gerror.New("插件未安装，请安装后操作！")
	}

	// 执行插件结构迁移（建表），幂等
	if err = migrateUp(m.Ctx(), m.GetSkeleton().Name); err != nil {
		return
	}

	// 执行插件数据填充（种子），幂等
	if err = migrateSeed(m.Ctx(), m.GetSkeleton().Name); err != nil {
		return
	}

	data := g.Map{
		dao.SysAddonsInstall.Columns().Version: m.GetSkeleton().Version,
	}
	return g.DB().Transaction(m.Ctx(), func(ctx context.Context, tx gdb.TX) error {
		if _, err = GetModel(ctx).Where(dao.SysAddonsInstall.Columns().Id, record.Id).Data(data).Update(); err != nil {
			return err
		}
		return m.Upgrade(ctx)
	})
}

// UnInstall 卸载模块
func UnInstall(m Module) (err error) {
	record, err := ScanInstall(m)
	if err != nil {
		return
	}

	if record == nil || record.Status != consts.AddonsInstallStatusOk {
		return gerror.New("插件未安装，请安装后操作！")
	}

	data := g.Map{
		dao.SysAddonsInstall.Columns().Version: m.GetSkeleton().Version,
		dao.SysAddonsInstall.Columns().Status:  consts.AddonsInstallStatusUn,
	}
	return g.DB().Transaction(m.Ctx(), func(ctx context.Context, tx gdb.TX) error {
		if _, err = GetModel(ctx).Where(dao.SysAddonsInstall.Columns().Id, record.Id).Data(data).Update(); err != nil {
			return err
		}
		return m.UnInstall(ctx)
	})
}

// isDevelopMode 是否为开发模式（只有 develop 才允许磁盘优先，product 使用内嵌 SQL）
func isDevelopMode(ctx context.Context) bool {
	return g.Cfg().MustGet(ctx, "system.mode", "develop").String() == "develop"
}

// migrateUp 执行插件结构迁移（建表），SQL 使用 CREATE TABLE IF NOT EXISTS 保证幂等，无需版本号。
// 开发模式下优先读磁盘，改 SQL 不用重新编译；生产模式优先读二进制内嵌 SQL。
func migrateUp(ctx context.Context, name string) (err error) {
	dir := gfile.Join("addons", name, "migrations")
	if isDevelopMode(ctx) && gfile.Exists(dir) {
		return migrate.Seed(ctx, dir)
	}

	sub := path.Join(name, "migrations")
	if _, statErr := fs.Stat(addonsFS.SQLFS, sub); statErr == nil {
		return migrate.SeedFS(ctx, addonsFS.SQLFS, sub)
	}
	if !gfile.Exists(dir) {
		return nil
	}
	return migrate.Seed(ctx, dir)
}

// migrateSeed 执行插件数据填充（种子），SQL 幂等
func migrateSeed(ctx context.Context, name string) (err error) {
	dir := gfile.Join("addons", name, "seeders")
	if isDevelopMode(ctx) && gfile.Exists(dir) {
		return migrate.Seed(ctx, dir)
	}

	sub := path.Join(name, "seeders")
	if _, statErr := fs.Stat(addonsFS.SQLFS, sub); statErr == nil {
		return migrate.SeedFS(ctx, addonsFS.SQLFS, sub)
	}
	if !gfile.Exists(dir) {
		return nil
	}
	return migrate.Seed(ctx, dir)
}
