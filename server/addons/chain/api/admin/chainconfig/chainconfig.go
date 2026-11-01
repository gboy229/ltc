// Package chainconfig
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package chainconfig

import (
	"github.com/gogf/gf/v2/frame/g"
	"hotgo/addons/chain/model/input/sysin"
	"hotgo/internal/model/input/form"
)

// ListReq 查询网络配置列表
type ListReq struct {
	g.Meta `path:"/chainconfig/list" method:"get" tags:"网络配置" summary:"获取网络配置列表"`
	sysin.ChainConfigListInp
}

type ListRes struct {
	List []*sysin.ChainConfigListModel `json:"list" dc:"数据列表"`
	form.PageRes
}

// ViewReq 获取指定信息
type ViewReq struct {
	g.Meta `path:"/chainconfig/view" method:"get" tags:"网络配置" summary:"获取指定网络配置信息"`
	sysin.ChainConfigViewInp
}

type ViewRes struct {
	*sysin.ChainConfigViewModel
}

// EditReq 修改/新增
type EditReq struct {
	g.Meta `path:"/chainconfig/edit" method:"post" tags:"网络配置" summary:"修改/新增网络配置"`
	sysin.ChainConfigEditInp
}

type EditRes struct{}

// DeleteReq 删除
type DeleteReq struct {
	g.Meta `path:"/chainconfig/delete" method:"post" tags:"网络配置" summary:"删除网络配置"`
	sysin.ChainConfigDeleteInp
}

type DeleteRes struct{}

// GenerateAddressReq 生成提币地址
type GenerateAddressReq struct {
	g.Meta `path:"/chainconfig/generateAddress" method:"post" tags:"网络配置" summary:"生成提币地址"`
	sysin.ChainConfigGenerateAddressInp
}

type GenerateAddressRes struct {
	*sysin.ChainConfigGenerateAddressModel
}

// OptionsReq 获取网络配置选项
type OptionsReq struct {
	g.Meta `path:"/chainconfig/options" method:"get" tags:"网络配置" summary:"获取网络配置选项"`
}

type OptionsRes struct {
	List []*sysin.ChainConfigOptionModel `json:"list" dc:"链选项"`
}
