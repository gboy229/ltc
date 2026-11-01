// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import "context"

type (
	IUserWallet interface {
		// GetWallet 获取用户登录钱包地址
		GetWallet(ctx context.Context, userId int64) (string, error)
		// GetUserIdsByWallet 按钱包地址模糊查询用户ID
		GetUserIdsByWallet(ctx context.Context, wallet string) (ids []uint64, err error)
		// GetWalletsByUserIds 批量获取用户钱包地址
		GetWalletsByUserIds(ctx context.Context, userIds []uint64) (wallets map[uint64]string, err error)
		// VerifySignature 校验当前用户钱包签名并消费 nonce
		VerifySignature(ctx context.Context, userId int64, nonce, signature string) error
	}
)

type (
	IUserRelation interface {
		// GetDirectUpline 获取用户的一级上级
		GetDirectUpline(ctx context.Context, userId int64) (int64, error)
	}
)

var (
	localUserWallet   IUserWallet
	localUserRelation IUserRelation
)

func UserWallet() IUserWallet {
	if localUserWallet == nil {
		panic("implement not found for interface IUserWallet, forgot register?")
	}
	return localUserWallet
}

func RegisterUserWallet(i IUserWallet) {
	localUserWallet = i
}

func UserRelation() IUserRelation {
	if localUserRelation == nil {
		panic("implement not found for interface IUserRelation, forgot register?")
	}
	return localUserRelation
}

func RegisterUserRelation(i IUserRelation) {
	localUserRelation = i
}
