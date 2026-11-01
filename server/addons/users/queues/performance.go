// Package queues
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package queues

import (
	"context"
	"encoding/json"

	"hotgo/addons/users/service"
	"hotgo/internal/library/queue"
)

const performanceTopic = "users_performance"

type performanceEvent struct {
	UserId uint64 `json:"userId"`
}

func init() {
	queue.RegisterConsumer(Performance)
}

// Performance 用户业绩刷新消费者
var Performance = &qPerformance{}

type qPerformance struct{}

func (q *qPerformance) GetTopic() string {
	return performanceTopic
}

func (q *qPerformance) Handle(ctx context.Context, mqMsg queue.MqMsg) error {
	var event performanceEvent
	if err := json.Unmarshal(mqMsg.Body, &event); err != nil {
		return err
	}
	return service.Performance().RefreshByUser(ctx, int64(event.UserId))
}
