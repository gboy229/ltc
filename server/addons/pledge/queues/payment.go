// Package queues
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package queues

import (
	"context"
	"encoding/json"

	iconsts "hotgo/internal/consts"
	"hotgo/internal/library/queue"
	mainmodel "hotgo/internal/model"

	"hotgo/addons/pledge/service"
)

func init() {
	queue.RegisterConsumer(Payment)
}

// Payment 链上支付确认事件消费者
var Payment = &qPayment{}

type qPayment struct{}

func (q *qPayment) GetTopic() string {
	return iconsts.QueueChainPayment
}

func (q *qPayment) Handle(ctx context.Context, mqMsg queue.MqMsg) (err error) {
	var event mainmodel.ChainPaymentEvent
	if err = json.Unmarshal(mqMsg.Body, &event); err != nil {
		return
	}
	return service.Order().HandlePayment(ctx, &event)
}
