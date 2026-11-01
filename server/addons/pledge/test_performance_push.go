//go:build ignore

// 业绩刷新消息测试命令，测试完成后可直接删除本文件。
//
// 在 server 目录执行：
//
//	go run addons/pledge/test_performance_push.go -user-id 123
//
// 只调用 perf.Push(ctx, userId)，不创建、不修改任何业务数据。
package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	perf "hotgo/addons/pledge/logic/performance"
)

func main() {
	userId := flag.Uint64("user-id", 0, "要推送业绩刷新的真实用户ID")
	flag.Parse()

	if *userId == 0 {
		fmt.Println("请传真实的 -user-id")
		flag.Usage()
		os.Exit(1)
	}

	if err := perf.Push(context.Background(), *userId); err != nil {
		fmt.Println("推送业绩刷新消息失败:", err)
		os.Exit(1)
	}
	fmt.Printf("已推送业绩刷新消息 userId=%d topic=%s\n", *userId, perf.Topic)
}
