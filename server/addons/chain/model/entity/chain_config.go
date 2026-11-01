// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainConfig is the golang structure for table chain_config.
type ChainConfig struct {
	Id              uint64      `json:"id"              orm:"id"               description:""`
	Driver          string      `json:"driver"          orm:"driver"           description:"驱动"`
	Chain           string      `json:"chain"           orm:"chain"            description:"网络"`
	RpcUrl          string      `json:"rpcUrl"          orm:"rpc_url"          description:"rpc地址"`
	ChainId         int         `json:"chainId"         orm:"chain_id"         description:"链id"`
	BrowserUrl      string      `json:"browserUrl"      orm:"browser_url"      description:"浏览器地址"`
	TransferAddress string      `json:"transferAddress" orm:"transfer_address" description:"转账地址"`
	TransferPrivate string      `json:"transferPrivate" orm:"transfer_private" description:"转账私钥"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:""`
	Confirmations   int         `json:"confirmations"   orm:"confirmations"    description:"确认数"`
}
