// Package user
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/library/walletsign"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/internal/service"
)

type sUserWallet struct{}

func NewUserWallet() *sUserWallet {
	return &sUserWallet{}
}

func init() {
	service.RegisterUserWallet(NewUserWallet())
}

// GetWallet 获取用户登录钱包地址
func (s *sUserWallet) GetWallet(ctx context.Context, userId int64) (string, error) {
	if userId <= 0 {
		return "", gerror.New("请先登录")
	}

	var account *entity.UsersUserAccount
	if err := dao.UsersUserAccount.Ctx(ctx).
		Where(dao.UsersUserAccount.Columns().UserId, userId).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		Scan(&account); err != nil {
		return "", gerror.Wrap(err, "获取钱包地址失败")
	}
	if account == nil || account.Account == "" {
		return "", gerror.New("用户未绑定钱包地址")
	}
	return account.Account, nil
}

// GetUserIdsByWallet 按钱包地址模糊查询用户ID
func (s *sUserWallet) GetUserIdsByWallet(ctx context.Context, wallet string) (ids []uint64, err error) {
	if wallet == "" {
		return []uint64{}, nil
	}

	var accounts []*entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		WhereLike(dao.UsersUserAccount.Columns().Account, "%"+wallet+"%").
		Scan(&accounts); err != nil {
		return nil, gerror.Wrap(err, "按钱包地址查询用户失败")
	}

	ids = make([]uint64, 0, len(accounts))
	for _, account := range accounts {
		ids = append(ids, uint64(account.UserId))
	}
	return
}

// GetWalletsByUserIds 批量获取用户钱包地址
func (s *sUserWallet) GetWalletsByUserIds(ctx context.Context, userIds []uint64) (wallets map[uint64]string, err error) {
	wallets = make(map[uint64]string, len(userIds))
	if len(userIds) == 0 {
		return
	}

	var accounts []*entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		WhereIn(dao.UsersUserAccount.Columns().UserId, userIds).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		Scan(&accounts); err != nil {
		return nil, gerror.Wrap(err, "批量查询钱包地址失败")
	}

	for _, account := range accounts {
		wallets[uint64(account.UserId)] = account.Account
	}
	return
}

// VerifySignature 校验当前用户钱包签名并消费 nonce
func (s *sUserWallet) VerifySignature(ctx context.Context, userId int64, nonce, signature string) error {
	address, err := s.GetWallet(ctx, userId)
	if err != nil {
		return err
	}
	return walletsign.VerifyAndConsume(ctx, address, nonce, signature)
}
