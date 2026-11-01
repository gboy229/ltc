// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"github.com/gogf/gf/v2/os/gtime"
)

// ChainAddress is the golang structure for table chain_address.
type ChainAddress struct {
	Id        uint64      `json:"id"        orm:"id"         description:""`
	Chain     string      `json:"chain"     orm:"chain"      description:"链名"`
	Uid       int         `json:"uid"       orm:"uid"        description:"用户ID"`
	Address   string      `json:"address"   orm:"address"    description:"地址"`
	Keys      string      `json:"keys"      orm:"keys"       description:"私钥"`
	CreatedAt *gtime.Time `json:"createdAt" orm:"created_at" description:""`
	UpdatedAt *gtime.Time `json:"updatedAt" orm:"updated_at" description:""`
}
