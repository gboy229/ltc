// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"

	"hotgo/addons/users/api/api/team"
)

type (
	ITeam interface {
		// Push 我的直推
		Push(ctx context.Context, userId int64, in *team.PushReq) (list []*team.PushItem, totalCount int, err error)
	}
)

var (
	localTeam ITeam
)

func Team() ITeam {
	if localTeam == nil {
		panic("implement not found for interface ITeam, forgot register?")
	}
	return localTeam
}

func RegisterTeam(i ITeam) {
	localTeam = i
}
