// Package asset
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package asset

import (
	"github.com/gogf/gf/v2/frame/g"

	"hotgo/addons/users/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// AccountListReq 资产账户列表
type AccountListReq struct {
	g.Meta `path:"/asset/accountList" method:"get" tags:"资产账户" summary:"获取资产账户列表"`
	sysin.AssetAccountListInp
}

type AccountListRes struct {
	List []*sysin.AssetAccountListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// LogListReq 资金流水列表
type LogListReq struct {
	g.Meta `path:"/asset/logList" method:"get" tags:"资金流水" summary:"获取资金流水列表"`
	sysin.AssetLogListInp
}

type LogListRes struct {
	List []*sysin.AssetLogListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// StatReq 资产统计
type StatReq struct {
	g.Meta `path:"/asset/stat" method:"get" tags:"资产统计" summary:"获取资产统计"`
	sysin.AssetStatInp
}

type StatRes struct {
	*sysin.AssetStatModel
}
