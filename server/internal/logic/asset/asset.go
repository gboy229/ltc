// Package asset
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package asset

import (
	"context"
	"strconv"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/shopspring/decimal"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/assetin"
	"hotgo/internal/service"
	"hotgo/utility/validate"
)

type sAsset struct{}

func NewAsset() *sAsset {
	return &sAsset{}
}

func init() {
	service.RegisterAsset(NewAsset())
}

// accountRow 行锁读取的账户，金额使用字符串避免 decimal 转 float64 丢精度
type accountRow struct {
	Id        uint64
	Available string
	Frozen    string
	Version   int
	Status    int
}

// EnsureAccount 确保用户资产账户存在
func (s *sAsset) EnsureAccount(ctx context.Context, userId uint64, account string) (err error) {
	if userId == 0 {
		return gerror.New("用户ID不能为空")
	}
	if account == "" {
		account = consts.AssetAccountUSDT
	}

	cols := dao.SysAssetAccount.Columns()
	_, err = dao.SysAssetAccount.Ctx(ctx).Data(g.Map{
		cols.UserId:    userId,
		cols.Account:   account,
		cols.Available: "0.00000000",
		cols.Frozen:    "0.00000000",
		cols.Status:    consts.StatusEnabled,
		cols.CreatedAt: gtime.Now(),
		cols.UpdatedAt: gtime.Now(),
	}).InsertIgnore()
	if err != nil {
		return gerror.Wrap(err, "创建资产账户失败")
	}
	return
}

// Change 资产变更（增减/冻结/解冻/结算），同一业务同一动作幂等
func (s *sAsset) Change(ctx context.Context, in *assetin.ChangeInp) (err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}

	num, err := parseAmount(strconv.FormatFloat(in.Num, 'f', -1, 64))
	if err != nil {
		err = gerror.New("变动金额格式不正确")
		return
	}

	// 重放快路径：已处理过的业务直接返回成功，不进入事务、不拿行锁
	if processed, checkErr := s.isProcessed(ctx, in); checkErr != nil {
		err = checkErr
		return
	} else if processed {
		return nil
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		colsAcc := dao.SysAssetAccount.Columns()
		colsLog := dao.SysAssetLog.Columns()

		// 确保账户存在，重复插入直接忽略
		if _, err = dao.SysAssetAccount.Ctx(ctx).
			Data(g.Map{
				colsAcc.UserId:    in.UserId,
				colsAcc.Account:   in.Account,
				colsAcc.Available: "0.00000000",
				colsAcc.Frozen:    "0.00000000",
				colsAcc.Status:    consts.StatusEnabled,
				colsAcc.CreatedAt: gtime.Now(),
				colsAcc.UpdatedAt: gtime.Now(),
			}).
			InsertIgnore(); err != nil {
			return gerror.Wrap(err, "创建资产账户失败")
		}

		// 锁定账户行，同一用户的资产变更串行执行
		acc := new(accountRow)
		if err = dao.SysAssetAccount.Ctx(ctx).
			Where(colsAcc.UserId, in.UserId).
			Where(colsAcc.Account, in.Account).
			LockUpdate().
			Scan(acc); err != nil {
			return gerror.Wrap(err, "锁定资产账户失败")
		}
		if acc.Id == 0 {
			return gerror.New("资产账户不存在")
		}

		// 锁内二次校验，防止并发重复写入
		processed, err := s.isProcessed(ctx, in)
		if err != nil {
			return gerror.Wrap(err, "查询资产流水失败")
		}
		if processed {
			return nil
		}

		beforeAvail, err := parseAmount(acc.Available)
		if err != nil {
			return gerror.Wrap(err, "解析可用余额失败")
		}
		beforeFrozen, err := parseAmount(acc.Frozen)
		if err != nil {
			return gerror.Wrap(err, "解析冻结金额失败")
		}

		afterAvail, afterFrozen, deltaAvail, deltaFrozen, err := calcChange(in.ChangeType, num, beforeAvail, beforeFrozen)
		if err != nil {
			return
		}

		result, err := dao.SysAssetAccount.Ctx(ctx).
			Where(colsAcc.Id, acc.Id).
			Where(colsAcc.Version, acc.Version).
			Data(g.Map{
				colsAcc.Available: afterAvail.StringFixed(8),
				colsAcc.Frozen:    afterFrozen.StringFixed(8),
				colsAcc.Version:   acc.Version + 1,
				colsAcc.UpdatedAt: gtime.Now(),
			}).
			Update()
		if err != nil {
			return gerror.Wrap(err, "更新资产账户失败")
		}
		if affected, _ := result.RowsAffected(); affected == 0 {
			return gerror.New("资产账户更新冲突，请重试")
		}

		if _, err = dao.SysAssetLog.Ctx(ctx).
			Data(g.Map{
				colsLog.UserId:          in.UserId,
				colsLog.Account:         in.Account,
				colsLog.AddonName:       in.AddonName,
				colsLog.ChangeType:      in.ChangeType,
				colsLog.BizType:         in.BizType,
				colsLog.BizId:           in.BizId,
				colsLog.BizSn:           in.BizSn,
				colsLog.BeforeAvailable: beforeAvail.StringFixed(8),
				colsLog.ChangeAvailable: deltaAvail.StringFixed(8),
				colsLog.AfterAvailable:  afterAvail.StringFixed(8),
				colsLog.BeforeFrozen:    beforeFrozen.StringFixed(8),
				colsLog.ChangeFrozen:    deltaFrozen.StringFixed(8),
				colsLog.AfterFrozen:     afterFrozen.StringFixed(8),
				colsLog.Remark:          in.Remark,
				colsLog.OperatorId:      in.OperatorId,
				colsLog.CreatedAt:       gtime.Now(),
			}).
			InsertAndGetId(); err != nil {
			return gerror.Wrap(err, "写入资产流水失败")
		}

		return nil
	})
	return
}

// isProcessed 判断同一业务、同一账务动作是否已处理
func (s *sAsset) isProcessed(ctx context.Context, in *assetin.ChangeInp) (bool, error) {
	colsLog := dao.SysAssetLog.Columns()
	return dao.SysAssetLog.Ctx(ctx).
		Where(colsLog.Account, in.Account).
		Where(colsLog.AddonName, in.AddonName).
		Where(colsLog.BizType, in.BizType).
		Where(colsLog.BizId, in.BizId).
		Where(colsLog.ChangeType, in.ChangeType).
		Exist()
}

// Account 获取用户资产账户
func (s *sAsset) Account(ctx context.Context, in *assetin.AccountInp) (res *assetin.AccountModel, err error) {
	if err = validate.PreFilter(ctx, in); err != nil {
		return
	}
	if in.Account == "" {
		in.Account = consts.AssetAccountUSDT
	}

	var row *entity.SysAssetAccount
	if err = dao.SysAssetAccount.Ctx(ctx).
		Where(dao.SysAssetAccount.Columns().UserId, in.UserId).
		Where(dao.SysAssetAccount.Columns().Account, in.Account).
		Scan(&row); err != nil {
		return
	}

	res = new(assetin.AccountModel)
	if row != nil {
		res.SysAssetAccount = *row
	}
	res.UserId = in.UserId
	return
}

// AccountList 获取资产账户列表
func (s *sAsset) AccountList(ctx context.Context, in *assetin.AccountListInp) (list []*assetin.AccountListModel, totalCount int, err error) {
	cols := dao.SysAssetAccount.Columns()
	mod := dao.SysAssetAccount.Ctx(ctx)
	if in.UserId > 0 {
		mod = mod.Where(cols.UserId, in.UserId)
	}
	if len(in.UserIds) > 0 {
		mod = mod.WhereIn(cols.UserId, in.UserIds)
	}
	if in.Account != "" {
		mod = mod.Where(cols.Account, in.Account)
	}

	if totalCount, err = mod.Count(); err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return
	}
	return
}

// Logs 获取资产流水列表
func (s *sAsset) Logs(ctx context.Context, in *assetin.LogListInp) (list []*assetin.LogListModel, totalCount int, err error) {
	cols := dao.SysAssetLog.Columns()
	mod := dao.SysAssetLog.Ctx(ctx)

	if in.UserId > 0 {
		mod = mod.Where(cols.UserId, in.UserId)
	}
	if len(in.UserIds) > 0 {
		mod = mod.WhereIn(cols.UserId, in.UserIds)
	}
	if in.Account != "" {
		mod = mod.Where(cols.Account, in.Account)
	}
	if in.AddonName != "" {
		mod = mod.Where(cols.AddonName, in.AddonName)
	}
	if in.ChangeType != "" {
		mod = mod.Where(cols.ChangeType, in.ChangeType)
	}
	if in.Direction == "income" {
		mod = mod.WhereGT(cols.ChangeAvailable, 0)
	} else if in.Direction == "expense" {
		mod = mod.WhereLT(cols.ChangeAvailable, 0)
	}
	if in.BizType != "" {
		mod = mod.Where(cols.BizType, in.BizType)
	}
	if in.BizId > 0 {
		mod = mod.Where(cols.BizId, in.BizId)
	}
	if len(in.CreatedAt) == 2 && in.CreatedAt[0] != nil && in.CreatedAt[1] != nil {
		mod = mod.WhereBetween(cols.CreatedAt, in.CreatedAt[0], in.CreatedAt[1])
	}

	if totalCount, err = mod.Count(); err != nil {
		return
	}
	if totalCount == 0 {
		return
	}

	if err = mod.Page(in.Page, in.PerPage).OrderDesc(cols.Id).Scan(&list); err != nil {
		return
	}
	return
}

// Stat 资产统计
func (s *sAsset) Stat(ctx context.Context, in *assetin.StatInp) (res *assetin.StatModel, err error) {
	res = new(assetin.StatModel)
	colsAcc := dao.SysAssetAccount.Columns()
	accMod := dao.SysAssetAccount.Ctx(ctx)
	if in.UserId > 0 {
		accMod = accMod.Where(colsAcc.UserId, in.UserId)
	}

	if res.AccountCount, err = accMod.Count(); err != nil {
		return
	}

	if res.TotalAvailable, err = accMod.Sum(colsAcc.Available); err != nil {
		return
	}
	if res.TotalFrozen, err = accMod.Sum(colsAcc.Frozen); err != nil {
		return
	}
	return
}

// calcChange 计算变动后的可用/冻结金额
func calcChange(changeType string, num, beforeAvail, beforeFrozen decimal.Decimal) (afterAvail, afterFrozen, deltaAvail, deltaFrozen decimal.Decimal, err error) {
	afterAvail, afterFrozen = beforeAvail, beforeFrozen

	switch changeType {
	case consts.AssetChangeTypeIncr:
		afterAvail = beforeAvail.Add(num)
		deltaAvail = num
	case consts.AssetChangeTypeDecr:
		if beforeAvail.LessThan(num) {
			err = gerror.New("余额不足")
			return
		}
		afterAvail = beforeAvail.Sub(num)
		deltaAvail = num.Neg()
	case consts.AssetChangeTypeFreeze:
		if beforeAvail.LessThan(num) {
			err = gerror.New("可用余额不足，无法冻结")
			return
		}
		afterAvail = beforeAvail.Sub(num)
		afterFrozen = beforeFrozen.Add(num)
		deltaAvail = num.Neg()
		deltaFrozen = num
	case consts.AssetChangeTypeUnfreeze:
		if beforeFrozen.LessThan(num) {
			err = gerror.New("冻结金额不足，无法解冻")
			return
		}
		afterAvail = beforeAvail.Add(num)
		afterFrozen = beforeFrozen.Sub(num)
		deltaAvail = num
		deltaFrozen = num.Neg()
	case consts.AssetChangeTypeSettle:
		if beforeFrozen.LessThan(num) {
			err = gerror.New("冻结金额不足，无法结算")
			return
		}
		afterFrozen = beforeFrozen.Sub(num)
		deltaFrozen = num.Neg()
	default:
		err = gerror.New("账务动作不合法")
	}
	return
}

// parseAmount 解析金额字符串
func parseAmount(s string) (decimal.Decimal, error) {
	if s == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(s)
}
