// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

// UsersTree is the golang structure for table users_tree.
type UsersTree struct {
	Id    uint64 `json:"id"    orm:"id"    description:""`
	Uid   int    `json:"uid"   orm:"uid"   description:"uid"`
	Pid   int    `json:"pid"   orm:"pid"   description:"pid"`
	Level int    `json:"level" orm:"level" description:"level"`
}
