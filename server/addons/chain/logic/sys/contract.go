// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"hotgo/addons/chain/logic/adapters"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/addons/chain/service"
)

type sSysContract struct{}

func NewSysContract() *sSysContract {
	return &sSysContract{}
}

func init() {
	service.RegisterSysContract(NewSysContract())
}

// Delete 删除合约，支持批量
func (s *sSysContract) Delete(ctx context.Context, in *sysin.ContractDeleteInp) (err error) {
	cols := dao.ChainContract.Columns()

	count, err := dao.ChainContract.Ctx(ctx).WhereIn(cols.Id, in.Ids).Count()
	if err != nil {
		return
	}
	if count == 0 {
		err = gerror.New("合约不存在或已删除")
		return
	}

	_, err = dao.ChainContract.Ctx(ctx).WhereIn(cols.Id, in.Ids).Delete()
	return
}

// Edit 修改/新增合约
func (s *sSysContract) Edit(ctx context.Context, in *sysin.ContractEditInp) (err error) {
	cols := dao.ChainContract.Columns()

	// 修改：允许调整所有参数；链名/合约地址变更时重新校验链配置、地址格式与唯一性
	if in.Id > 0 {
		var current *entity.ChainContract
		if err = dao.ChainContract.Ctx(ctx).WherePri(in.Id).Scan(&current); err != nil {
			return
		}
		if current == nil {
			err = gerror.New("合约不存在或已删除")
			return
		}

		if current.Chain != in.Chain || !strings.EqualFold(current.Address, in.Address) {
			configCols := dao.ChainConfig.Columns()
			var config *entity.ChainConfig
			if err = dao.ChainConfig.Ctx(ctx).Where(configCols.Chain, in.Chain).Scan(&config); err != nil {
				return
			}
			if config == nil {
				err = gerror.Newf("链配置不存在：%s，请先在网络配置中添加该链", in.Chain)
				return
			}

			adapter, adapterErr := adapters.NewChainAdapter(config.Driver)
			if adapterErr != nil {
				return gerror.Wrap(adapterErr, "不支持的链驱动，无法校验合约地址")
			}
			if !adapter.IsValidAddress(in.Address) {
				err = gerror.Newf("合约地址格式不正确，请检查 %s 链（%s）的地址格式", in.Chain, config.Driver)
				return
			}

			addressCount, countErr := dao.ChainContract.Ctx(ctx).
				Where(cols.Chain, in.Chain).
				Where("LOWER("+cols.Address+") = ?", strings.ToLower(in.Address)).
				WhereNot(cols.Id, in.Id).
				Count()
			if countErr != nil {
				return countErr
			}
			if addressCount > 0 {
				err = gerror.New("该链下已存在相同合约地址")
				return
			}
		}

		_, err = dao.ChainContract.Ctx(ctx).Fields(sysin.ContractUpdateFields{}).WherePri(in.Id).Data(in).Update()
		return
	}

	// 新增：链配置必须存在，并取链驱动用于地址校验
	configCols := dao.ChainConfig.Columns()
	var config *entity.ChainConfig
	if err = dao.ChainConfig.Ctx(ctx).Where(configCols.Chain, in.Chain).Scan(&config); err != nil {
		return
	}
	if config == nil {
		err = gerror.Newf("链配置不存在：%s，请先在网络配置中添加该链", in.Chain)
		return
	}

	// 地址格式由链驱动校验，EVM、Tron 等各自实现，不在业务层写死
	adapter, err := adapters.NewChainAdapter(config.Driver)
	if err != nil {
		err = gerror.Wrap(err, "不支持的链驱动，无法校验合约地址")
		return
	}
	if !adapter.IsValidAddress(in.Address) {
		err = gerror.Newf("合约地址格式不正确，请检查 %s 链（%s）的地址格式", in.Chain, config.Driver)
		return
	}

	// 同一链下合约地址唯一（地址大小写不敏感）
	mod := dao.ChainContract.Ctx(ctx).
		Where(cols.Chain, in.Chain).
		Where("LOWER("+cols.Address+") = ?", strings.ToLower(in.Address))
	addressCount, err := mod.Count()
	if err != nil {
		return
	}
	if addressCount > 0 {
		err = gerror.New("该链下已存在相同合约地址")
		return
	}

	_, err = dao.ChainContract.Ctx(ctx).Fields(sysin.ContractInsertFields{}).Data(in).Insert()
	return
}

// View 获取指定合约信息
func (s *sSysContract) View(ctx context.Context, in *sysin.ContractViewInp) (res *sysin.ContractViewModel, err error) {
	if err = dao.ChainContract.Ctx(ctx).WherePri(in.Id).Scan(&res); err != nil {
		return
	}
	return
}

// List 获取合约列表
func (s *sSysContract) List(ctx context.Context, in *sysin.ContractListInp) (list []*sysin.ContractListModel, totalCount int, err error) {
	mod := dao.ChainContract.Ctx(ctx)
	cols := dao.ChainContract.Columns()

	if in.Chain != "" {
		mod = mod.WhereLike(cols.Chain, "%"+in.Chain+"%")
	}

	if in.Symbol != "" {
		mod = mod.WhereLike(cols.Symbol, "%"+in.Symbol+"%")
	}

	if in.Address != "" {
		mod = mod.WhereLike(cols.Address, "%"+in.Address+"%")
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
