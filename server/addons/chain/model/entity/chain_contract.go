// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainContract is the golang structure for table chain_contract.
type ChainContract struct {
	Id             uint64      `json:"id"             orm:"id"              description:""`
	Chain          string      `json:"chain"          orm:"chain"           description:"链名"`
	Address        string      `json:"address"        orm:"address"         description:"合约地址"`
	Symbol         string      `json:"symbol"         orm:"symbol"          description:"币种"`
	Decimal        int         `json:"decimal"        orm:"decimal"         description:"精度"`
	MinAmount      float64     `json:"minAmount"      orm:"min_amount"      description:""`
	MaxAmount      float64     `json:"maxAmount"      orm:"max_amount"      description:""`
	MultipleAmount int         `json:"multipleAmount" orm:"multiple_amount" description:"倍数"`
	VerifyAmount   float64     `json:"verifyAmount"   orm:"verify_amount"   description:""`
	Scale          int         `json:"scale"          orm:"scale"           description:"保留几位小数"`
	DayLimit       int         `json:"dayLimit"       orm:"day_limit"       description:"日上限笔数"`
	FeeType        int         `json:"feeType"        orm:"fee_type"        description:"手续费类型，0=百分比，1=个数"`
	Rate           float64     `json:"rate"           orm:"rate"            description:""`
	RechargeSwitch int         `json:"rechargeSwitch" orm:"recharge_switch" description:"充币开关"`
	WithdrawSwitch int         `json:"withdrawSwitch" orm:"withdraw_switch" description:"提现开关"`
	CreatedAt      *gtime.Time `json:"createdAt"      orm:"created_at"      description:""`
	UpdatedAt      *gtime.Time `json:"updatedAt"      orm:"updated_at"      description:""`
	ScanBlock      uint64      `json:"scanBlock"      orm:"scan_block"      description:"已扫块高"`
	ScanBlockHash  string      `json:"scanBlockHash"  orm:"scan_block_hash" description:"已扫块哈希"`
}
