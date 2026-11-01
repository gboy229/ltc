// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/logic/adapters"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/addons/chain/service"
	"hotgo/addons/chain/utility"
)

type sSysChainConfig struct{}

func NewSysChainConfig() *sSysChainConfig {
	return &sSysChainConfig{}
}

func init() {
	service.RegisterSysChainConfig(NewSysChainConfig())
}

// Delete 删除网络配置，支持批量
func (s *sSysChainConfig) Delete(ctx context.Context, in *sysin.ChainConfigDeleteInp) (err error) {
	cols := dao.ChainConfig.Columns()

	count, err := dao.ChainConfig.Ctx(ctx).WhereIn(cols.Id, in.Ids).Count()
	if err != nil {
		return
	}
	if count == 0 {
		err = gerror.New("网络配置不存在或已删除")
		return
	}

	_, err = dao.ChainConfig.Ctx(ctx).WhereIn(cols.Id, in.Ids).Delete()
	return
}

// Edit 修改/新增网络配置
func (s *sSysChainConfig) Edit(ctx context.Context, in *sysin.ChainConfigEditInp) (err error) {
	cols := dao.ChainConfig.Columns()

	// 同一驱动下链名唯一
	mod := dao.ChainConfig.Ctx(ctx).Where(cols.Driver, in.Driver).Where(cols.Chain, in.Chain)
	if in.Id > 0 {
		mod = mod.WhereNot(cols.Id, in.Id)
	}
	count, err := mod.Count()
	if err != nil {
		return
	}
	if count > 0 {
		err = gerror.New("该驱动下已存在同名网络配置")
		return
	}

	// 修改
	if in.Id > 0 {
		_, err = dao.ChainConfig.Ctx(ctx).Fields(sysin.ChainConfigUpdateFields{}).WherePri(in.Id).Data(in).Update()
		return
	}

	// 新增
	_, err = dao.ChainConfig.Ctx(ctx).Fields(sysin.ChainConfigInsertFields{}).Data(in).OmitEmptyData().Insert()
	return
}

// View 获取指定网络配置信息
func (s *sSysChainConfig) View(ctx context.Context, in *sysin.ChainConfigViewInp) (res *sysin.ChainConfigViewModel, err error) {
	if err = dao.ChainConfig.Ctx(ctx).WherePri(in.Id).Scan(&res); err != nil {
		return
	}
	return
}

// List 获取网络配置列表
func (s *sSysChainConfig) List(ctx context.Context, in *sysin.ChainConfigListInp) (list []*sysin.ChainConfigListModel, totalCount int, err error) {
	mod := dao.ChainConfig.Ctx(ctx)
	cols := dao.ChainConfig.Columns()

	if in.Driver != "" {
		mod = mod.WhereLike(cols.Driver, "%"+in.Driver+"%")
	}

	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}

	totalCount, err = mod.Count()
	if err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).Order("id desc").Scan(&list); err != nil {
		return
	}
	return
}

// GenerateAddress 生成提币地址，私钥加密后落库，不对外返回
func (s *sSysChainConfig) GenerateAddress(ctx context.Context, in *sysin.ChainConfigGenerateAddressInp) (res *sysin.ChainConfigGenerateAddressModel, err error) {
	configCols := dao.ChainConfig.Columns()

	var models *entity.ChainConfig
	// 只查询生成地址所需的字段，避免把历史私钥读入内存
	if err = dao.ChainConfig.Ctx(ctx).Fields(configCols.Id, configCols.Driver, configCols.Chain).WherePri(in.Id).Scan(&models); err != nil {
		return
	}
	if models == nil {
		err = gerror.New("网络配置不存在或已删除")
		return
	}

	// 通过链驱动生成地址
	adapter, err := adapters.NewChainAdapter(models.Driver)
	if err != nil {
		err = gerror.Wrap(err, "不支持的链驱动")
		return
	}
	adapter.Chain(models.Chain)

	address, privateKey, err := adapter.GenerateAddress()
	if err != nil {
		err = gerror.Wrap(err, "生成钱包地址失败")
		return
	}
	if address == "" || privateKey == "" {
		err = gerror.New("生成钱包地址失败")
		return
	}

	// 私钥加密后落库，禁止明文存储
	encryptedPrivateKey, err := utility.EncryptPrivateKey(ctx, privateKey)
	if err != nil {
		return
	}

	err = dao.ChainConfig.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		addressCols := dao.ChainAddress.Columns()

		// 保存地址历史记录
		count, err := dao.ChainAddress.Ctx(ctx).Where(addressCols.Address, address).Count()
		if err != nil {
			return gerror.Wrap(err, "检查钱包地址失败")
		}
		if count > 0 {
			return gerror.New("钱包地址已存在，请重试")
		}

		if _, err = dao.ChainAddress.Ctx(ctx).Data(g.Map{
			addressCols.Chain:   models.Chain,
			addressCols.Uid:     0,
			addressCols.Address: address,
			addressCols.Keys:    encryptedPrivateKey,
		}).Insert(); err != nil {
			return gerror.Wrap(err, "保存钱包地址失败")
		}

		// 更新当前提币地址，私钥同样以密文保存
		if _, err = dao.ChainConfig.Ctx(ctx).WherePri(models.Id).Data(g.Map{
			configCols.TransferAddress: address,
			configCols.TransferPrivate: encryptedPrivateKey,
		}).Update(); err != nil {
			return gerror.Wrap(err, "更新提币地址失败")
		}
		return nil
	})
	if err != nil {
		return
	}

	res = &sysin.ChainConfigGenerateAddressModel{
		Address: address,
	}
	return
}

// Options 获取网络配置选项
func (s *sSysChainConfig) Options(ctx context.Context) (list []*sysin.ChainConfigOptionModel, err error) {
	cols := dao.ChainConfig.Columns()

	var models []*entity.ChainConfig
	if err = dao.ChainConfig.Ctx(ctx).Fields(cols.Id, cols.Chain).Order("id asc").Scan(&models); err != nil {
		return
	}

	list = make([]*sysin.ChainConfigOptionModel, 0, len(models))
	for _, m := range models {
		list = append(list, &sysin.ChainConfigOptionModel{
			Label: m.Chain,
			Value: m.Chain,
		})
	}
	return
}
