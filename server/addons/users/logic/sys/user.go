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
	"github.com/gogf/gf/v2/os/gtime"
	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/model/input/sysin"
	"hotgo/addons/users/service"
	"hotgo/addons/users/utility"
)

type sSysUser struct{}

func NewSysUser() *sSysUser {
	return &sSysUser{}
}

func init() {
	service.RegisterSysUser(NewSysUser())
}

// Delete 删除用户（连同登录凭证）
func (s *sSysUser) Delete(ctx context.Context, in *sysin.UserDeleteInp) (err error) {
	var models *entity.UsersUser
	if err = dao.UsersUser.Ctx(ctx).WherePri(in.Id).Scan(&models); err != nil {
		return
	}
	if models == nil {
		err = gerror.New("用户不存在或已删除")
		return
	}

	return dao.UsersUser.Transaction(ctx, func(ctx context.Context, tx gdb.TX) (err error) {
		if _, err = dao.UsersUser.Ctx(ctx).WherePri(in.Id).Delete(); err != nil {
			return
		}
		_, err = dao.UsersUserAccount.Ctx(ctx).Where(dao.UsersUserAccount.Columns().UserId, in.Id).Delete()
		return
	})
}

// Edit 修改/新增用户
func (s *sSysUser) Edit(ctx context.Context, in *sysin.UserEditInp) (res *sysin.UserEditModel, err error) {
	res = new(sysin.UserEditModel)

	// 修改
	if in.Id > 0 {
		if _, err = dao.UsersUser.Ctx(ctx).Fields(sysin.UserUpdateFields{}).WherePri(in.Id).Data(in).Update(); err != nil {
			return
		}
		res.Id = in.Id
		return
	}

	// 新增：自动生成对外唯一用户号
	if in.Uid == "" {
		in.Uid = utility.GenUid()
	}
	id, err := dao.UsersUser.Ctx(ctx).Fields(sysin.UserInsertFields{}).Data(in).OmitEmptyData().InsertAndGetId()
	if err != nil {
		return
	}
	in.Id = id
	res.Id = id
	return
}

// Status 更新用户状态
func (s *sSysUser) Status(ctx context.Context, in *sysin.UserStatusInp) (err error) {
	_, err = dao.UsersUser.Ctx(ctx).WherePri(in.Id).Data(dao.UsersUser.Columns().Status, in.Status).Update()
	return
}

// View 获取指定用户信息
func (s *sSysUser) View(ctx context.Context, in *sysin.UserViewInp) (res *sysin.UserViewModel, err error) {
	if err = dao.UsersUser.Ctx(ctx).WherePri(in.Id).Scan(&res); err != nil {
		return
	}
	if res == nil {
		return
	}

	// 等级名称
	if res.LevelId > 0 {
		levelName, _ := dao.UsersLevel.Ctx(ctx).Fields(dao.UsersLevel.Columns().Name).WherePri(res.LevelId).Value()
		res.LevelName = levelName.String()
	}

	var performance *entity.UsersPerformance
	if err = dao.UsersPerformance.Ctx(ctx).WherePri(res.Id).Scan(&performance); err != nil {
		return
	}
	if performance != nil {
		res.PersonalAmount = performance.PersonalAmount
		res.DirectAmount = performance.DirectAmount
		res.TeamAmount = performance.TeamAmount
	}

	return
}

// List 获取用户列表
func (s *sSysUser) List(ctx context.Context, in *sysin.UserListInp) (list []*sysin.UserListModel, totalCount int, err error) {
	mod := dao.UsersUser.Ctx(ctx)
	cols := dao.UsersUser.Columns()

	if in.Uid != "" {
		mod = mod.WhereLike(cols.Uid, "%"+in.Uid+"%")
	}
	if in.Nickname != "" {
		mod = mod.WhereLike(cols.Nickname, "%"+in.Nickname+"%")
	}
	if in.LevelId > 0 {
		mod = mod.Where(cols.LevelId, in.LevelId)
	}
	if in.Status > 0 {
		mod = mod.Where(cols.Status, in.Status)
	}
	if in.AuthStatus >= 0 {
		mod = mod.Where(cols.AuthStatus, in.AuthStatus)
	}
	if len(in.CreatedAt) == 2 {
		mod = mod.WhereBetween(cols.CreatedAt, gtime.New(in.CreatedAt[0]), gtime.New(in.CreatedAt[1]))
	}
	// 钱包地址搜索
	if in.WalletAddress != "" {
		ids, e := s.resolveUserIdsByWallet(ctx, in.WalletAddress)
		if e != nil {
			err = e
			return
		}
		if len(ids) == 0 {
			return
		}
		mod = mod.WhereIn(cols.Id, ids)
	}
	// 推荐人钱包搜索
	if in.ParentWallet != "" {
		ids, e := s.resolveUserIdsByWallet(ctx, in.ParentWallet)
		if e != nil {
			err = e
			return
		}
		if len(ids) == 0 {
			return
		}
		mod = mod.WhereIn(cols.Pid, ids)
	}
	// 团队钱包搜索：先定位用户，再查其全部下级
	if in.TeamWallet != "" {
		ids, e := s.resolveUserIdsByWallet(ctx, in.TeamWallet)
		if e != nil {
			err = e
			return
		}
		if len(ids) == 0 {
			return
		}
		var treeRows []*entity.UsersTree
		if err = dao.UsersTree.Ctx(ctx).
			Fields(dao.UsersTree.Columns().Uid).
			WhereIn(dao.UsersTree.Columns().Pid, ids).
			Scan(&treeRows); err != nil {
			return
		}
		descendantIds := make([]int64, 0, len(treeRows))
		seen := make(map[int64]struct{}, len(treeRows))
		for _, t := range treeRows {
			d := int64(t.Uid)
			if _, ok := seen[d]; !ok {
				seen[d] = struct{}{}
				descendantIds = append(descendantIds, d)
			}
		}
		if len(descendantIds) == 0 {
			return
		}
		mod = mod.WhereIn(cols.Id, descendantIds)
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

	err = s.fillListRelations(ctx, list)
	return
}

// fillListRelations 批量填充等级名称、推荐人信息、用户名
func (s *sSysUser) fillListRelations(ctx context.Context, list []*sysin.UserListModel) (err error) {
	if len(list) == 0 {
		return
	}

	// 等级名称
	levelIds := make([]int, 0)
	for _, v := range list {
		if v.LevelId > 0 {
			levelIds = append(levelIds, v.LevelId)
		}
	}
	if len(levelIds) > 0 {
		var levels []*entity.UsersLevel
		if err = dao.UsersLevel.Ctx(ctx).
			Fields(dao.UsersLevel.Columns().Id, dao.UsersLevel.Columns().Name).
			WhereIn(dao.UsersLevel.Columns().Id, levelIds).
			Scan(&levels); err != nil {
			return
		}
		levelMap := make(map[int]string, len(levels))
		for _, lv := range levels {
			levelMap[int(lv.Id)] = lv.Name
		}
		for _, v := range list {
			v.LevelName = levelMap[v.LevelId]
		}
	}

	// 推荐人信息
	pidIds := make([]int64, 0)
	for _, v := range list {
		if v.Pid > 0 {
			pidIds = append(pidIds, v.Pid)
		}
	}
	if len(pidIds) > 0 {
		var parents []*entity.UsersUser
		if err = dao.UsersUser.Ctx(ctx).
			Fields(dao.UsersUser.Columns().Id, dao.UsersUser.Columns().Uid, dao.UsersUser.Columns().Nickname).
			WhereIn(dao.UsersUser.Columns().Id, pidIds).
			Scan(&parents); err != nil {
			return
		}
		parentMap := make(map[int64]*entity.UsersUser, len(parents))
		for _, p := range parents {
			parentMap[p.Id] = p
		}
		for _, v := range list {
			if p, ok := parentMap[v.Pid]; ok {
				v.ParentUid = p.Uid
				v.ParentNickname = p.Nickname
			}
		}
	}

	// 钱包地址（wallet 登录凭证）
	userIds := make([]int64, 0, len(list))
	for _, v := range list {
		userIds = append(userIds, v.Id)
	}
	var accounts []*entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		Fields(dao.UsersUserAccount.Columns().UserId, dao.UsersUserAccount.Columns().Account).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		WhereIn(dao.UsersUserAccount.Columns().UserId, userIds).
		Scan(&accounts); err != nil {
		return
	}
	walletMap := make(map[int64]string, len(accounts))
	for _, a := range accounts {
		if _, ok := walletMap[a.UserId]; !ok {
			walletMap[a.UserId] = a.Account
		}
	}
	for _, v := range list {
		v.WalletAddress = walletMap[v.Id]
	}

	// 业绩
	var performances []*entity.UsersPerformance
	if err = dao.UsersPerformance.Ctx(ctx).
		WhereIn(dao.UsersPerformance.Columns().UserId, userIds).
		Scan(&performances); err != nil {
		return
	}
	performanceMap := make(map[int64]*entity.UsersPerformance, len(performances))
	for _, p := range performances {
		performanceMap[int64(p.UserId)] = p
	}
	for _, v := range list {
		if p, ok := performanceMap[v.Id]; ok {
			v.PersonalAmount = p.PersonalAmount
			v.DirectAmount = p.DirectAmount
			v.TeamAmount = p.TeamAmount
		}
	}
	return
}

// resolveUserIdsByWallet 根据钱包地址模糊匹配用户ID列表
func (s *sSysUser) resolveUserIdsByWallet(ctx context.Context, wallet string) (ids []int64, err error) {
	if wallet == "" {
		return nil, nil
	}
	accountCols := dao.UsersUserAccount.Columns()
	var accounts []*entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		Fields(accountCols.UserId).
		Where(accountCols.Type, usersconsts.UserAuthTypeWallet).
		WhereLike(accountCols.Account, "%"+wallet+"%").
		Scan(&accounts); err != nil {
		return nil, gerror.Wrap(err, "查询钱包用户失败")
	}
	ids = make([]int64, 0, len(accounts))
	seen := make(map[int64]struct{}, len(accounts))
	for _, a := range accounts {
		if _, ok := seen[a.UserId]; !ok {
			seen[a.UserId] = struct{}{}
			ids = append(ids, a.UserId)
		}
	}
	return ids, nil
}
