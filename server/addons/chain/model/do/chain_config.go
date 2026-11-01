// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package do

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainConfig is the golang structure of table hg_addon_chain_config for DAO operations like Where/Data.
type ChainConfig struct {
	g.Meta          `orm:"table:hg_addon_chain_config, do:true"`
	Id              interface{} //
	Driver          interface{} // 驱动
	Chain           interface{} // 网络
	RpcUrl          interface{} // rpc地址
	ChainId         interface{} // 链id
	BrowserUrl      interface{} // 浏览器地址
	TransferAddress interface{} // 转账地址
	TransferPrivate interface{} // 转账私钥
	CreatedAt       *gtime.Time //
	UpdatedAt       *gtime.Time //
	Confirmations   interface{} // 确认数
}
