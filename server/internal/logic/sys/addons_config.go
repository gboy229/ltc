// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2023 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/library/contexts"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"
	"hotgo/utility/simple"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gstr"
)

var AddonsMaskDemoField []string

type sSysAddonsConfig struct{}

func NewSysAddonsConfig() *sSysAddonsConfig {
	return &sSysAddonsConfig{}
}

func init() {
	service.RegisterSysAddonsConfig(NewSysAddonsConfig())
}

// GetConfigByGroup 获取指定分组的配置
func (s *sSysAddonsConfig) GetConfigByGroup(ctx context.Context, in *sysin.GetAddonsConfigInp) (res *sysin.GetAddonsConfigModel, err error) {
	if in.AddonName == "" {
		err = gerror.New("插件名称不能为空")
		return
	}

	if in.Group == "" {
		err = gerror.New("分组不能为空")
		return
	}

	var (
		mod    = dao.SysAddonsConfig.Ctx(ctx)
		models []*entity.SysAddonsConfig
		cols   = dao.SysAddonsConfig.Columns()
	)

	if err = mod.Fields(cols.Key, cols.Value, cols.Type).
		Where(cols.AddonName, in.AddonName).
		Where(cols.Group, in.Group).
		Scan(&models); err != nil {
		return nil, err
	}

	res = new(sysin.GetAddonsConfigModel)
	res.List = make(g.Map)
	if len(models) > 0 {
		for _, v := range models {
			val, err := s.ConversionType(ctx, v)
			if err != nil {
				return nil, err
			}
			res.List[v.Key] = val
			if simple.IsDemo(ctx) && gstr.InArray(AddonsMaskDemoField, v.Key) {
				res.List[v.Key] = consts.DemoTips
			}
		}
	}
	return
}

// GetConfigByAddon 按插件名称获取配置 key-value，Redis 优先
func (s *sSysAddonsConfig) GetConfigByAddon(ctx context.Context, addonName string) (config map[string]string, err error) {
	if addonName == "" {
		return nil, gerror.New("插件名称不能为空")
	}

	cacheKey := consts.CacheAddonsConfig + addonName
	if value, cacheErr := g.Redis().Do(ctx, "HGETALL", cacheKey); cacheErr == nil && value != nil {
		if cached := value.MapStrStr(); len(cached) > 0 {
			return cached, nil
		}
	}

	cols := dao.SysAddonsConfig.Columns()
	var models []*entity.SysAddonsConfig
	if err = dao.SysAddonsConfig.Ctx(ctx).
		Where(cols.AddonName, addonName).
		Scan(&models); err != nil {
		return nil, gerror.Wrap(err, "查询插件配置失败")
	}

	config = make(map[string]string, len(models))
	for _, model := range models {
		config[model.Key] = model.Value
	}
	if len(config) == 0 {
		return config, nil
	}

	args := make([]interface{}, 0, len(config)*2+1)
	args = append(args, cacheKey)
	for k, v := range config {
		args = append(args, k, v)
	}
	if _, cacheErr := g.Redis().Do(ctx, "HSET", args...); cacheErr != nil {
		g.Log().Warningf(ctx, "写入插件配置缓存失败 addon:%s err:%+v", addonName, cacheErr)
	} else {
		_, _ = g.Redis().Do(ctx, "EXPIRE", cacheKey, 3600)
	}
	return config, nil
}

// GetConfigByKey 按插件名称 + key 获取配置值，Redis 优先
func (s *sSysAddonsConfig) GetConfigByKey(ctx context.Context, addonName, key string) (value string, err error) {
	if addonName == "" {
		return "", gerror.New("插件名称不能为空")
	}
	if key == "" {
		return "", gerror.New("配置key不能为空")
	}

	cacheKey := consts.CacheAddonsConfig + addonName
	if cached, cacheErr := g.Redis().Do(ctx, "HGET", cacheKey, key); cacheErr == nil && cached != nil && !cached.IsNil() && !cached.IsEmpty() {
		return cached.String(), nil
	}

	config, err := s.GetConfigByAddon(ctx, addonName)
	if err != nil {
		return "", err
	}
	return config[key], nil
}

// GetCurrentAddonConfig 按当前请求插件获取全部配置
func (s *sSysAddonsConfig) GetCurrentAddonConfig(ctx context.Context) (config map[string]string, err error) {
	addonName := contexts.GetAddonName(ctx)
	if addonName == "" {
		return nil, gerror.New("当前请求不属于插件，无法获取插件配置")
	}
	return s.GetConfigByAddon(ctx, addonName)
}

// GetCurrentAddonConfigByKey 按当前请求插件 + key 获取配置值
func (s *sSysAddonsConfig) GetCurrentAddonConfigByKey(ctx context.Context, key string) (value string, err error) {
	addonName := contexts.GetAddonName(ctx)
	if addonName == "" {
		return "", gerror.New("当前请求不属于插件，无法获取插件配置")
	}
	return s.GetConfigByKey(ctx, addonName, key)
}

// GetConfigListByAddon 获取插件全部配置项（含元数据）
func (s *sSysAddonsConfig) GetConfigListByAddon(ctx context.Context, addonName string) (list []*entity.SysAddonsConfig, err error) {
	if addonName == "" {
		err = gerror.New("插件名称不能为空")
		return
	}

	cols := dao.SysAddonsConfig.Columns()
	if err = dao.SysAddonsConfig.Ctx(ctx).
		Where(cols.AddonName, addonName).
		OrderAsc(cols.Sort).
		Scan(&list); err != nil {
		return nil, err
	}
	return
}

// ConversionType 转换类型
func (s *sSysAddonsConfig) ConversionType(ctx context.Context, models *entity.SysAddonsConfig) (value interface{}, err error) {
	if models == nil {
		err = gerror.New("数据不存在")
		return
	}
	return consts.ConvType(models.Value, models.Type), nil
}

// UpdateConfigByGroup 更新指定分组的配置
func (s *sSysAddonsConfig) UpdateConfigByGroup(ctx context.Context, in *sysin.UpdateAddonsConfigInp) (err error) {
	if in.AddonName == "" {
		err = gerror.New("插件名称不能为空")
		return
	}

	if in.Group == "" {
		err = gerror.New("分组不能为空")
		return
	}

	var (
		mod    = dao.SysAddonsConfig.Ctx(ctx)
		models []*entity.SysAddonsConfig
	)

	if err = mod.
		Where("addon_name", in.AddonName).
		Where("group", in.Group).
		Scan(&models); err != nil {
		return
	}

	err = dao.SysAddonsConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		for k, v := range in.List {
			row := s.getConfigByKey(k, models)
			// 新增
			if row == nil {
				return gerror.Newf("暂不支持从前台添加变量，请先在数据库表[%v]中配置变量：%v", dao.SysAddonsConfig.Table(), k)
			}

			// 更新
			if _, err = dao.SysAddonsConfig.Ctx(ctx).Where("id", row.Id).Data(g.Map{"value": v, "updated_at": gtime.Now()}).Update(); err != nil {
				err = gerror.Wrap(err, consts.ErrorORM)
				return
			}
		}

		return
	})
	if err == nil {
		_, _ = g.Redis().Do(ctx, "DEL", consts.CacheAddonsConfig+in.AddonName)
	}
	return
}

func (s *sSysAddonsConfig) getConfigByKey(key string, models []*entity.SysAddonsConfig) *entity.SysAddonsConfig {
	if len(models) == 0 {
		return nil
	}

	for _, v := range models {
		if key == v.Key {
			return v
		}
	}
	return nil
}
