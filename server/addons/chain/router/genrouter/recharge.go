// Package genrouter
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package genrouter

import "hotgo/addons/chain/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.Recharge) // 充值记录
}
