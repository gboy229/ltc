// Package genrouter
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package genrouter

import "hotgo/addons/pledge/controller/admin/sys"

func init() {
	LoginRequiredRouter = append(LoginRequiredRouter, sys.Plan) // 投资套餐
}
