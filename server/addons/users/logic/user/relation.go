// Package user
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package user

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"

	"hotgo/addons/users/model/dao"
	"hotgo/addons/users/model/entity"
	"hotgo/internal/service"
)

type sUserRelation struct{}

func NewUserRelation() *sUserRelation { return &sUserRelation{} }

func init() { service.RegisterUserRelation(NewUserRelation()) }

// GetDirectUpline 获取用户的一级上级
func (s *sUserRelation) GetDirectUpline(ctx context.Context, userId int64) (int64, error) {
	if userId <= 0 {
		return 0, nil
	}

	cols := dao.UsersTree.Columns()
	var tree *entity.UsersTree
	if err := dao.UsersTree.Ctx(ctx).
		Where(cols.Uid, userId).
		Where(cols.Level, 1).
		Scan(&tree); err != nil {
		return 0, gerror.Wrap(err, "获取直推上级失败")
	}
	if tree == nil {
		return 0, nil
	}
	return int64(tree.Pid), nil
}
