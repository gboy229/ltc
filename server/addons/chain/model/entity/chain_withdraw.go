// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainWithdraw is the golang structure for table chain_withdraw.
type ChainWithdraw struct {
	Id              uint64      `json:"id"              orm:"id"               description:""`
	Uid             int         `json:"uid"             orm:"uid"              description:"用户ID"`
	Chain           string      `json:"chain"           orm:"chain"            description:"链"`
	Symbol          string      `json:"symbol"          orm:"symbol"           description:"币种"`
	ContractAddress string      `json:"contractAddress" orm:"contract_address" description:"合约地址"`
	FromAddress     string      `json:"fromAddress"     orm:"from_address"     description:"转出地址"`
	ToAddress       string      `json:"toAddress"       orm:"to_address"       description:"转入地址"`
	PayAmount       float64     `json:"payAmount"       orm:"pay_amount"       description:"付款数量，含手续费"`
	ReceiveAmount   float64     `json:"receiveAmount"   orm:"receive_amount"   description:"收款数量"`
	Fee             float64     `json:"fee"             orm:"fee"              description:"手续费"`
	Rate            float64     `json:"rate"            orm:"rate"             description:"手续费费率"`
	Hash            string      `json:"hash"            orm:"hash"             description:"哈希"`
	Status          int         `json:"status"          orm:"status"           description:"状态，0=审核中，1=已审核，2=已打币，3=失败，4=完成，-1=取消"`
	Version         int         `json:"version"         orm:"version"          description:"乐观锁"`
	AuditBy         int64       `json:"auditBy"         orm:"audit_by"         description:"审核人ID"`
	AuditAt         *gtime.Time `json:"auditAt"         orm:"audit_at"         description:"审核时间"`
	Remark          string      `json:"remark"          orm:"remark"           description:"审核备注"`
	CreatedAt       *gtime.Time `json:"createdAt"       orm:"created_at"       description:""`
	UpdatedAt       *gtime.Time `json:"updatedAt"       orm:"updated_at"       description:""`
}
