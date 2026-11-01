//go:build ignore

// 直推奖励测试命令，测试完成后可直接删除本文件。
//
// 不创建模拟数据，直接使用数据库中的真实订单、上级关系和周期配置。
//
// 在 server 目录执行：
//
//	go run addons/pledge/test_push_reward.go -order-id 123 -amount 3 -period-no 1
//
// 参数：
//
//	-order-id   下级真实订单 ID
//	-amount     本次静态收益金额，通常等于 order.amount * 月化率 / 100 / 30
//	-period-no  收益期数，不传默认 order.inc_income_num + 1
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/v2/frame/g"

	_ "hotgo/addons/pledge/logic/income"
	_ "hotgo/addons/pledge/logic/settlement"
	"hotgo/addons/pledge/service"
	_ "hotgo/addons/users/logic/user"
	_ "hotgo/internal/logic/asset"
)

func main() {
	orderId := flag.Uint64("order-id", 0, "下级真实订单ID")
	amount := flag.Float64("amount", 0, "本次静态收益金额")
	periodNo := flag.Int("period-no", 0, "收益期数，0表示自动取订单已结算期数+1")
	flag.Parse()

	if *orderId == 0 || *amount <= 0 {
		fmt.Println("请传真实的 -order-id 和 -amount")
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()
	if err := service.Settlement().CreatePushReward(ctx, *orderId, *amount, *periodNo); err != nil {
		fmt.Println("生成直推奖励失败:", err)
		os.Exit(1)
	}

	fmt.Printf("直推奖励处理完成 orderId=%d amount=%.8f periodNo=%d\n", *orderId, *amount, *periodNo)

	var rows []struct {
		UserId   uint64  `json:"user_id"`
		Type     string  `json:"type"`
		Amount   float64 `json:"amount"`
		Rate     float64 `json:"rate"`
		OrderId  uint64  `json:"order_id"`
		FromUser uint64  `json:"from_user_id"`
	}
	if err := g.DB().Model("hg_addon_pledge_income_log").
		Fields("user_id,type,amount,rate,order_id,from_user_id").
		Where("order_id", *orderId).
		Where("type", "direct").
		OrderDesc("id").
		Scan(&rows); err != nil {
		fmt.Println("查询直推收益失败:", err)
		os.Exit(1)
	}

	fmt.Println("直推收益流水：")
	for _, row := range rows {
		fmt.Printf("  user=%d amount=%.8f rate=%.4f%% order=%d from=%d\n", row.UserId, row.Amount, row.Rate, row.OrderId, row.FromUser)
	}
}
