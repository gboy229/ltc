// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"hotgo/addons/users/api/api/center"
)

type (
	IUserCenter interface {
		// Info 当前用户信息
		Info(ctx context.Context, userId int64, in *center.InfoReq) (res *center.InfoRes, err error)
	}
)

var (
	localUserCenter IUserCenter
)

func UserCenter() IUserCenter {
	if localUserCenter == nil {
		panic("implement not found for interface IUserCenter, forgot register?")
	}
	return localUserCenter
}

func RegisterUserCenter(i IUserCenter) {
	localUserCenter = i
}
