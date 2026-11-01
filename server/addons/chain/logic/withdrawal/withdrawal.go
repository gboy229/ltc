// Package withdrawal
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package withdrawal

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	apiWithdrawal "hotgo/addons/chain/api/api/withdrawal"
	chainconsts "hotgo/addons/chain/consts"
	"hotgo/addons/chain/model/dao"
	"hotgo/addons/chain/model/entity"
	"hotgo/addons/chain/service"
	"hotgo/internal/consts"
	"hotgo/internal/library/ratelimit"
	"hotgo/internal/model/input/assetin"
	basesvc "hotgo/internal/service"
)

type sWithdrawal struct{}

func NewWithdrawal() *sWithdrawal {
	return &sWithdrawal{}
}

func init() {
	service.RegisterWithdrawal(NewWithdrawal())
}

// Options 获取当前可用的提币网络
func (s *sWithdrawal) Options(ctx context.Context, in *apiWithdrawal.OptionsReq) (res *apiWithdrawal.OptionsRes, err error) {
	contractCols := dao.ChainContract.Columns()
	var contracts []*entity.ChainContract
	if err = dao.ChainContract.Ctx(ctx).
		Where(contractCols.WithdrawSwitch, chainconsts.ChainSwitchOn).
		OrderAsc(contractCols.Id).
		Scan(&contracts); err != nil {
		return nil, gerror.Wrap(err, "查询提币网络失败")
	}

	res = new(apiWithdrawal.OptionsRes)
	res.List = make([]*apiWithdrawal.Option, 0, len(contracts))
	for _, item := range contracts {
		res.List = append(res.List, &apiWithdrawal.Option{
			Chain:     item.Chain,
			Network:   networkLabel(item.Chain),
			Symbol:    strings.ToUpper(item.Symbol),
			Decimals:  item.Decimal,
			FeeType:   item.FeeType,
			Rate:      item.Rate,
			MinAmount: item.MinAmount,
			MaxAmount: item.MaxAmount,
		})
	}
	return
}

// Apply 提交提币申请并冻结资产
func (s *sWithdrawal) Apply(ctx context.Context, userId int64, wallet string, in *apiWithdrawal.ApplyReq) (res *apiWithdrawal.ApplyRes, err error) {

	if err = basesvc.UserWallet().VerifySignature(ctx, userId, in.Nonce, in.Signature); err != nil {
		return nil, err
	}

	allowed, err := ratelimit.Allow(ctx, fmt.Sprintf("rate:withdraw:apply:%d", userId), 10*time.Second)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, gerror.New("操作过于频繁，请稍后再试")
	}

	chain := strings.TrimSpace(in.Chain)
	symbol := strings.ToUpper(strings.TrimSpace(in.Symbol))
	toAddress := strings.TrimSpace(wallet)
	if chain == "" || symbol == "" {
		return nil, gerror.New("请选择提币网络")
	}
	if !common.IsHexAddress(toAddress) {
		return nil, gerror.New("提币地址格式不正确")
	}
	amount := int64(in.Amount)
	if amount <= 0 {
		return nil, gerror.New("提币数量必须大于0")
	}
	amountFloat := float64(amount)

	contract, err := s.getContract(ctx, chain, symbol)
	if err != nil {
		return nil, err
	}
	if contract == nil {
		return nil, gerror.New("该网络暂不支持提币")
	}
	if contract.WithdrawSwitch != chainconsts.ChainSwitchOn {
		return nil, gerror.New("该网络提币已关闭")
	}

	minAmount := contract.MinAmount
	if minAmount <= 0 {
		minAmount = 10
	}
	if amountFloat < minAmount {
		return nil, gerror.Newf("最少提币 %s %s", strconv.FormatFloat(minAmount, 'f', -1, 64), symbol)
	}
	if contract.MaxAmount > 0 && amountFloat > contract.MaxAmount {
		return nil, gerror.Newf("单笔最多提币 %s %s", strconv.FormatFloat(contract.MaxAmount, 'f', -1, 64), symbol)
	}

	fee := calcFee(amountFloat, contract.FeeType, contract.Rate)
	if fee < 0 || fee >= amountFloat {
		return nil, gerror.New("手续费配置异常")
	}
	receiveAmount := amountFloat - fee

	account, err := basesvc.Asset().Account(ctx, &assetin.AccountInp{
		UserId:  uint64(userId),
		Account: consts.AssetAccountUSDT,
	})
	if err != nil {
		return nil, err
	}
	if account == nil || account.Available < amountFloat {
		return nil, gerror.New("可用余额不足")
	}

	var chainConfig *entity.ChainConfig
	configCols := dao.ChainConfig.Columns()
	if err = dao.ChainConfig.Ctx(ctx).
		Where(configCols.Chain, chain).
		Scan(&chainConfig); err != nil {
		return nil, gerror.Wrap(err, "查询链配置失败")
	}
	if chainConfig == nil || chainConfig.TransferAddress == "" {
		return nil, gerror.New("提币地址未配置")
	}

	status := chainconsts.ChainWithdrawStatusAudited
	if amountFloat >= contract.VerifyAmount {
		status = chainconsts.ChainWithdrawStatusPending
	}

	cols := dao.ChainWithdraw.Columns()
	id, err := dao.ChainWithdraw.Ctx(ctx).Data(g.Map{
		cols.Uid:             userId,
		cols.Chain:           chain,
		cols.Symbol:          symbol,
		cols.ContractAddress: contract.Address,
		cols.FromAddress:     chainConfig.TransferAddress,
		cols.ToAddress:       toAddress,
		cols.PayAmount:       amount,
		cols.ReceiveAmount:   receiveAmount,
		cols.Fee:             fee,
		cols.Rate:            contract.Rate,
		cols.Status:          status,
		cols.Version:         0,
	}).InsertAndGetId()
	if err != nil {
		return nil, gerror.Wrap(err, "创建提币申请失败")
	}

	// 冻结用户可用资产，后台审核通过后再结算
	if err = basesvc.Asset().Change(ctx, &assetin.ChangeInp{
		UserId:     uint64(userId),
		Account:    consts.AssetAccountUSDT,
		AddonName:  consts.AssetAddonChain,
		ChangeType: consts.AssetChangeTypeFreeze,
		BizType:    consts.AssetBizTypeChainWithdraw,
		BizId:      uint64(id),
		Num:        amountFloat,
		Remark:     "提币申请冻结",
	}); err != nil {
		_, _ = dao.ChainWithdraw.Ctx(ctx).
			WherePri(id).
			Data(g.Map{cols.Status: chainconsts.ChainWithdrawStatusCanceled, cols.Remark: "冻结资产失败"}).
			Update()
		return nil, gerror.Wrap(err, "冻结资产失败")
	}

	return
}

// List 用户提币记录
func (s *sWithdrawal) List(ctx context.Context, userId int64, in *apiWithdrawal.ListReq) (res *apiWithdrawal.ListRes, err error) {

	cols := dao.ChainWithdraw.Columns()
	mod := dao.ChainWithdraw.Ctx(ctx).Where(cols.Uid, userId)

	switch strings.TrimSpace(in.Status) {
	case "processing":
		mod = mod.WhereIn(cols.Status, []int{
			chainconsts.ChainWithdrawStatusPending,
			chainconsts.ChainWithdrawStatusAudited,
			chainconsts.ChainWithdrawStatusTransferred,
		})
	case "success":
		mod = mod.WhereIn(cols.Status, []int{
			chainconsts.ChainWithdrawStatusFinished,
		})
	}

	totalCount, err := mod.Clone().Count()
	if err != nil {
		return nil, gerror.Wrap(err, "查询提币记录失败")
	}

	var records []*entity.ChainWithdraw
	if totalCount > 0 {
		if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&records); err != nil {
			return nil, gerror.Wrap(err, "查询提币记录失败")
		}
	}

	res = new(apiWithdrawal.ListRes)
	res.List = make([]*apiWithdrawal.Record, 0, len(records))
	for _, item := range records {
		date := ""
		if item.CreatedAt != nil {
			date = item.CreatedAt.Format("Y-m-d H:i")
		}
		res.List = append(res.List, &apiWithdrawal.Record{
			Chain:         item.Chain,
			Network:       networkLabel(item.Chain),
			Symbol:        strings.ToUpper(item.Symbol),
			Status:        userStatus(item.Status),
			Address:       item.ToAddress,
			Hash:          item.Hash,
			Date:          date,
			Amount:        item.PayAmount,
			Fee:           item.Fee,
			ReceiveAmount: item.ReceiveAmount,
		})
	}
	res.PageRes.Pack(in, totalCount)
	return
}

func (s *sWithdrawal) getContract(ctx context.Context, chain, symbol string) (*entity.ChainContract, error) {
	cols := dao.ChainContract.Columns()
	var contracts []*entity.ChainContract
	if err := dao.ChainContract.Ctx(ctx).
		Where(cols.Chain, chain).
		Scan(&contracts); err != nil {
		return nil, gerror.Wrap(err, "查询提币合约失败")
	}
	for _, item := range contracts {
		if strings.EqualFold(item.Symbol, symbol) {
			return item, nil
		}
	}
	return nil, nil
}

func calcFee(amount float64, feeType int, rate float64) float64 {
	if feeType == chainconsts.ChainFeeTypePercent {
		return round8(amount * rate / 100)
	}
	return round8(rate)
}

func round8(value float64) float64 {
	return math.Round(value*1e8) / 1e8
}

func networkLabel(chain string) string {
	switch strings.ToLower(chain) {
	case "bsc", "bnb":
		return "BEP20"
	case "eth", "ethereum":
		return "ERC20"
	case "tron", "trc20":
		return "TRC20"
	default:
		return strings.ToUpper(chain)
	}
}

func userStatus(status int) string {
	switch status {
	case chainconsts.ChainWithdrawStatusFinished:
		return "success"
	case chainconsts.ChainWithdrawStatusFailed, chainconsts.ChainWithdrawStatusCanceled:
		return "failed"
	default:
		return "processing"
	}
}

func withdrawId(item *entity.ChainWithdraw) string {
	if item.Hash != "" {
		return item.Hash
	}
	date := ""
	if item.CreatedAt != nil {
		date = item.CreatedAt.Format("YmdHis")
	}
	return fmt.Sprintf("%s-%s-%s", item.Chain, date, strconv.FormatFloat(item.PayAmount, 'f', -1, 64))
}
