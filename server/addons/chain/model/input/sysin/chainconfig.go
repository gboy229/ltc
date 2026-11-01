// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"hotgo/addons/chain/model/entity"
	"hotgo/internal/model/input/form"
)

// ChainConfigEditInp 修改/新增网络配置
type ChainConfigEditInp struct {
	entity.ChainConfig
}

func (in *ChainConfigEditInp) Filter(ctx context.Context) (err error) {
	if in.Driver == "" {
		err = gerror.New("驱动不能为空")
		return
	}
	if in.Chain == "" {
		err = gerror.New("网络不能为空")
		return
	}
	if in.ChainId <= 0 {
		err = gerror.New("链ID不能为空")
		return
	}
	if in.RpcUrl == "" {
		err = gerror.New("rpc地址不能为空")
		return
	}
	if in.BrowserUrl == "" {
		err = gerror.New("浏览器地址不能为空")
		return
	}
	// 确认数不在表单展示，未传时默认 15
	if in.Confirmations <= 0 {
		in.Confirmations = 15
	}
	if in.Confirmations > 1000 {
		err = gerror.New("确认数必须在0~1000之间")
		return
	}
	return
}

type ChainConfigEditModel struct{}

// ChainConfigUpdateFields 修改网络配置字段过滤
type ChainConfigUpdateFields struct {
	Driver        string `json:"driver"        description:"驱动"`
	Chain         string `json:"chain"         description:"网络"`
	RpcUrl        string `json:"rpcUrl"        description:"rpc地址"`
	ChainId       int    `json:"chainId"       description:"链id"`
	BrowserUrl    string `json:"browserUrl"    description:"浏览器地址"`
	Confirmations int    `json:"confirmations" description:"确认数"`
}

// ChainConfigInsertFields 新增网络配置字段过滤
type ChainConfigInsertFields struct {
	Driver        string `json:"driver"        description:"驱动"`
	Chain         string `json:"chain"         description:"网络"`
	RpcUrl        string `json:"rpcUrl"        description:"rpc地址"`
	ChainId       int    `json:"chainId"       description:"链id"`
	BrowserUrl    string `json:"browserUrl"    description:"浏览器地址"`
	Confirmations int    `json:"confirmations" description:"确认数"`
}

// ChainConfigDeleteInp 删除网络配置
type ChainConfigDeleteInp struct {
	Ids []int64 `json:"ids" v:"required#请选择要删除的网络配置" dc:"网络配置ID"`
}

func (in *ChainConfigDeleteInp) Filter(ctx context.Context) (err error) {
	if len(in.Ids) == 0 {
		err = gerror.New("请选择要删除的网络配置")
	}
	return
}

type ChainConfigDeleteModel struct{}

// ChainConfigViewInp 获取指定网络配置信息
type ChainConfigViewInp struct {
	Id int64 `json:"id" v:"required#网络配置ID不能为空" dc:"网络配置ID"`
}

// ChainConfigViewModel 网络配置详情，禁止返回转账私钥
type ChainConfigViewModel struct {
	Id              uint64      `json:"id"              orm:"id"              dc:"网络配置ID"`
	Driver          string      `json:"driver"          orm:"driver"          dc:"驱动"`
	Chain           string      `json:"chain"           orm:"chain"           dc:"网络"`
	RpcUrl          string      `json:"rpcUrl"          orm:"rpc_url"         dc:"rpc地址"`
	ChainId         int         `json:"chainId"         orm:"chain_id"        dc:"链id"`
	BrowserUrl      string      `json:"browserUrl"      orm:"browser_url"     dc:"浏览器地址"`
	Confirmations   int         `json:"confirmations"   orm:"confirmations"   dc:"确认数"`
	TransferAddress string      `json:"transferAddress" orm:"transfer_address" dc:"转账地址"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"      dc:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"      dc:""`
}

// ChainConfigListInp 获取网络配置列表
type ChainConfigListInp struct {
	form.PageReq
	Driver string `json:"driver" dc:"驱动"`
	Chain  string `json:"chain" dc:"网络"`
}

// ChainConfigListModel 网络配置列表，禁止返回转账私钥
type ChainConfigListModel struct {
	Id              uint64      `json:"id"              orm:"id"              dc:"网络配置ID"`
	Driver          string      `json:"driver"          orm:"driver"          dc:"驱动"`
	Chain           string      `json:"chain"           orm:"chain"           dc:"网络"`
	RpcUrl          string      `json:"rpcUrl"          orm:"rpc_url"         dc:"rpc地址"`
	ChainId         int         `json:"chainId"         orm:"chain_id"        dc:"链id"`
	BrowserUrl      string      `json:"browserUrl"      orm:"browser_url"     dc:"浏览器地址"`
	Confirmations   int         `json:"confirmations"   orm:"confirmations"   dc:"确认数"`
	TransferAddress string      `json:"transferAddress" orm:"transfer_address" dc:"转账地址"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"      dc:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"      dc:""`
}

// ChainConfigGenerateAddressInp 生成提币地址
type ChainConfigGenerateAddressInp struct {
	Id int64 `json:"id" v:"required#网络配置ID不能为空" dc:"网络配置ID"`
}

// ChainConfigGenerateAddressModel 生成提币地址结果，私钥不对外返回
type ChainConfigGenerateAddressModel struct {
	Address string `json:"address" dc:"转账地址"`
}

// ChainConfigOptionModel 网络配置下拉选项
type ChainConfigOptionModel struct {
	Label string `json:"label" dc:"显示名称"`
	Value string `json:"value" dc:"链名"`
}
