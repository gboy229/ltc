// Package performance
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package performance

import (
	"context"

	"hotgo/internal/library/queue"
)

const Topic = "users_performance"

type Event struct {
	UserId uint64 `json:"userId"`
}

// Push 推送用户业绩刷新事件
func Push(ctx context.Context, userId uint64) error {
	if userId == 0 {
		return nil
	}
	return queue.Push(Topic, &Event{UserId: userId})
}
