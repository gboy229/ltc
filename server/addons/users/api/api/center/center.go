// Package center
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package center

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/users/api/api/auth"
)

// InfoReq 当前用户信息
type InfoReq struct {
	g.Meta `path:"/center/info" method:"get" tags:"用户模块" summary:"当前用户信息"`
}

// InfoRes 当前用户信息响应
type InfoRes struct {
	User *auth.UserInfo `json:"user" dc:"用户信息"`
}
