// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainRecharge is the golang structure for table chain_recharge.
type ChainRecharge struct {
	Id              uint64      `json:"id"              orm:"id"               description:""`
	Type            string      `json:"type"            orm:"type"             description:""`
	Uid             int         `json:"uid"             orm:"uid"              description:"用户ID"`
	Chain           string      `json:"chain"           orm:"chain"            description:"链"`
	Symbol          string      `json:"symbol"          orm:"symbol"           description:"币种"`
	FromAddress     string      `json:"fromAddress"     orm:"from_address"     description:"转出地址"`
	ToAddress       string      `json:"toAddress"       orm:"to_address"       description:"转入地址"`
	ContractAddress string      `json:"contractAddress" orm:"contract_address" description:"合约地址"`
	Amount          float64     `json:"amount"          orm:"amount"           description:"充值数量"`
	BlockHeight     int64       `json:"blockHeight"     orm:"block_height"     description:"块高"`
	Hash            string      `json:"hash"            orm:"hash"             description:"交易哈希"`
	LogIndex        int         `json:"logIndex"        orm:"log_index"        description:"日志序号"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:""`
}
