// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"

	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/addons/users/service"
	basesvc "hotgo/internal/service"
)

type sSysAsset struct{}

func NewSysAsset() *sSysAsset {
	return &sSysAsset{}
}

func init() {
	service.RegisterSysAsset(NewSysAsset())
}

// AccountList 获取资产账户列表（补充用户昵称）
func (s *sSysAsset) AccountList(ctx context.Context, in *sysin.AssetAccountListInp) (list []*sysin.AssetAccountListModel, totalCount int, err error) {
	if in.Uid != "" || in.WalletAddress != "" {
		ids, resolveErr := s.resolveAccountUserIds(ctx, in.Uid, in.WalletAddress)
		if resolveErr != nil {
			return nil, 0, resolveErr
		}
		if len(ids) == 0 {
			return
		}
		in.UserIds = ids
	}

	accounts, totalCount, err := basesvc.Asset().AccountList(ctx, &in.AccountListInp)
	if err != nil || len(accounts) == 0 {
		return
	}

	userIds := make([]uint64, 0, len(accounts))
	for _, v := range accounts {
		userIds = append(userIds, v.UserId)
	}
	profiles, err := s.getUserProfiles(ctx, userIds)
	if err != nil {
		return
	}

	list = make([]*sysin.AssetAccountListModel, 0, len(accounts))
	for _, v := range accounts {
		profile := profiles[v.UserId]
		list = append(list, &sysin.AssetAccountListModel{
			AccountListModel: *v,
			Uid:              profile.Uid,
			Nickname:         profile.Nickname,
			WalletAddress:    profile.WalletAddress,
		})
	}
	return
}

// LogList 获取资金流水列表（补充用户昵称）
func (s *sSysAsset) LogList(ctx context.Context, in *sysin.AssetLogListInp) (list []*sysin.AssetLogListModel, totalCount int, err error) {
	if in.Uid != "" || in.WalletAddress != "" {
		ids, resolveErr := s.resolveAccountUserIds(ctx, in.Uid, in.WalletAddress)
		if resolveErr != nil {
			return nil, 0, resolveErr
		}
		if len(ids) == 0 {
			return
		}
		in.UserIds = ids
	}

	logs, totalCount, err := basesvc.Asset().Logs(ctx, &in.LogListInp)
	if err != nil || len(logs) == 0 {
		return
	}

	userIds := make([]uint64, 0, len(logs))
	for _, v := range logs {
		userIds = append(userIds, v.UserId)
	}
	profiles, err := s.getUserProfiles(ctx, userIds)
	if err != nil {
		return
	}

	list = make([]*sysin.AssetLogListModel, 0, len(logs))
	for _, v := range logs {
		profile := profiles[v.UserId]
		list = append(list, &sysin.AssetLogListModel{
			LogListModel:  *v,
			Uid:           profile.Uid,
			Nickname:      profile.Nickname,
			WalletAddress: profile.WalletAddress,
		})
	}
	return
}

// Stat 资产统计
func (s *sSysAsset) Stat(ctx context.Context, in *sysin.AssetStatInp) (res *sysin.AssetStatModel, err error) {
	data, err := basesvc.Asset().Stat(ctx, &in.StatInp)
	if err != nil {
		return
	}
	res = new(sysin.AssetStatModel)
	res.StatModel = *data
	return
}

type accountUserProfile struct {
	Uid           string
	Nickname      string
	WalletAddress string
}

// resolveAccountUserIds 按UID/钱包地址定位用户
func (s *sSysAsset) resolveAccountUserIds(ctx context.Context, uid, wallet string) ([]uint64, error) {
	userCols := dao.UsersUser.Columns()
	mod := dao.UsersUser.Ctx(ctx).Fields(dao.UsersUser.Table() + "." + userCols.Id)
	if uid != "" {
		mod = mod.WhereLike(userCols.Uid, "%"+uid+"%")
	}
	if wallet != "" {
		accountTable := dao.UsersUserAccount.Table()
		mod = mod.LeftJoin(accountTable+" a", "a.user_id="+dao.UsersUser.Table()+".id").
			Where("a.type", usersconsts.UserAuthTypeWallet).
			WhereLike("a.account", "%"+wallet+"%")
	}

	var rows []struct {
		Id uint64 `json:"id"`
	}
	if err := mod.Scan(&rows); err != nil {
		return nil, err
	}
	ids := make([]uint64, 0, len(rows))
	seen := make(map[uint64]struct{}, len(rows))
	for _, row := range rows {
		if _, ok := seen[row.Id]; ok {
			continue
		}
		seen[row.Id] = struct{}{}
		ids = append(ids, row.Id)
	}
	return ids, nil
}

// getUserProfiles 批量获取UID、昵称、钱包地址
func (s *sSysAsset) getUserProfiles(ctx context.Context, ids []uint64) (map[uint64]accountUserProfile, error) {
	userCols := dao.UsersUser.Columns()
	var users []struct {
		Id       uint64 `json:"id"`
		Uid      string `json:"uid"`
		Nickname string `json:"nickname"`
	}
	if err := dao.UsersUser.Ctx(ctx).
		Fields(userCols.Id, userCols.Uid, userCols.Nickname).
		WhereIn(userCols.Id, ids).
		Scan(&users); err != nil {
		return nil, err
	}

	profiles := make(map[uint64]accountUserProfile, len(users))
	for _, user := range users {
		profiles[user.Id] = accountUserProfile{Uid: user.Uid, Nickname: user.Nickname}
	}

	accountCols := dao.UsersUserAccount.Columns()
	var accounts []*entity.UsersUserAccount
	if err := dao.UsersUserAccount.Ctx(ctx).
		Fields(accountCols.UserId, accountCols.Account).
		Where(accountCols.Type, usersconsts.UserAuthTypeWallet).
		WhereIn(accountCols.UserId, ids).
		Scan(&accounts); err != nil {
		return nil, err
	}
	for _, account := range accounts {
		uid := uint64(account.UserId)
		profile := profiles[uid]
		if profile.WalletAddress == "" {
			profile.WalletAddress = account.Account
			profiles[uid] = profile
		}
	}
	return profiles, nil
}
