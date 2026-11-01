// Package user
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"

	apiAuth "hotgo/addons/users/api/api/auth"
	"hotgo/addons/users/api/api/center"
	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/library/walletcache"
	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/addons/users/service"
)

type sUserCenter struct{}

func NewUserCenter() *sUserCenter {
	return &sUserCenter{}
}

func init() {
	service.RegisterUserCenter(NewUserCenter())
}

// Info 当前用户信息
func (s *sUserCenter) Info(ctx context.Context, userId int64, in *center.InfoReq) (res *center.InfoRes, err error) {
	if userId <= 0 {
		err = gerror.New("用户身份异常，请重新登录")
		return
	}

	var user *entity.UsersUser
	if err = dao.UsersUser.Ctx(ctx).WherePri(userId).Scan(&user); err != nil {
		err = gerror.Wrap(err, "获取用户信息失败")
		return
	}
	if user == nil {
		err = gerror.New("用户不存在")
		return
	}

	// 钱包地址（wallet 登录凭证）
	var account *entity.UsersUserAccount
	if err = dao.UsersUserAccount.Ctx(ctx).
		Where(dao.UsersUserAccount.Columns().UserId, user.Id).
		Where(dao.UsersUserAccount.Columns().Type, usersconsts.UserAuthTypeWallet).
		Scan(&account); err != nil {
		err = gerror.Wrap(err, "获取钱包地址失败")
		return
	}

	// 会员等级
	var level *entity.UsersLevel
	if user.LevelId > 0 {
		if err = dao.UsersLevel.Ctx(ctx).WherePri(user.LevelId).Scan(&level); err != nil {
			err = gerror.Wrap(err, "获取会员等级失败")
			return
		}
	}

	address := ""
	if account != nil {
		address = account.Account
		// 兼容升级前已登录的用户：访问用户信息时补写 Redis 映射
		if cacheErr := walletcache.Set(ctx, address, user.Id); cacheErr != nil {
			g.Log().Warningf(ctx, "补写钱包用户缓存失败 uid:%d address:%s err:%+v", user.Id, address, cacheErr)
		}
	}
	levelName := ""
	if level != nil {
		levelName = level.Name
	}

	res = &center.InfoRes{User: &apiAuth.UserInfo{
		Uid:       user.Uid,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Address:   address,
		LevelName: levelName,
	}}
	return
}
