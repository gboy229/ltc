// Package assetin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package assetin

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"

	"hotgo/internal/consts"
	"hotgo/internal/library/dict"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/form"
)

// ChangeInp 资产变更
type ChangeInp struct {
	UserId     uint64  `json:"userId"     v:"required#用户ID不能为空" dc:"用户ID"`
	Account    string  `json:"account"    dc:"资产账户标识，默认USDT"`
	AddonName  string  `json:"addonName"  dc:"来源插件，默认system"`
	ChangeType string  `json:"changeType" v:"required#账务动作不能为空" dc:"账务动作"`
	BizType    string  `json:"bizType"    v:"required#业务类型不能为空" dc:"业务类型"`
	BizId      uint64  `json:"bizId"      v:"required#业务单号不能为空" dc:"业务单号"`
	BizSn      string  `json:"bizSn"      dc:"外部业务单号"`
	Num        float64 `json:"num"        v:"required#变动金额不能为空" dc:"变动金额（正数）"`
	Remark     string  `json:"remark"     dc:"备注"`
	OperatorId int64   `json:"operatorId" dc:"操作人ID"`
}

func (in *ChangeInp) Filter(ctx context.Context) (err error) {
	if in.UserId == 0 {
		err = gerror.New("用户ID不能为空")
		return
	}
	if in.Account == "" {
		in.Account = consts.AssetAccountUSDT
	}
	if in.AddonName == "" {
		in.AddonName = consts.AssetAddonSystem
	}
	if in.Num <= 0 {
		err = gerror.New("变动金额必须大于0")
		return
	}
	if !dict.HasOptionKey(consts.AssetChangeTypeOptions, in.ChangeType) {
		err = gerror.New("账务动作不合法")
		return
	}
	return
}

// AccountInp 获取资产账户
type AccountInp struct {
	UserId  uint64 `json:"userId"  v:"required#用户ID不能为空" dc:"用户ID"`
	Account string `json:"account" dc:"资产账户标识，默认USDT"`
}

// AccountModel 资产账户
type AccountModel struct {
	entity.SysAssetAccount
}

// AccountListInp 获取资产账户列表
type AccountListInp struct {
	form.PageReq
	UserId  uint64   `json:"userId"  dc:"用户ID"`
	UserIds []uint64 `json:"userIds" dc:"用户ID列表"`
	Account string   `json:"account" dc:"资产账户标识"`
}

// AccountListModel 资产账户
type AccountListModel struct {
	entity.SysAssetAccount
}

// LogListInp 获取资产流水列表
type LogListInp struct {
	form.PageReq
	UserId     uint64        `json:"userId"     dc:"用户ID"`
	UserIds    []uint64      `json:"userIds"    dc:"用户ID列表"`
	Account    string        `json:"account"    dc:"资产账户标识"`
	AddonName  string        `json:"addonName"  dc:"来源插件"`
	ChangeType string        `json:"changeType" dc:"账务动作"`
	Direction  string        `json:"direction"  dc:"资金方向：income/expense"`
	BizType    string        `json:"bizType"    dc:"业务类型"`
	BizId      uint64        `json:"bizId"      dc:"业务单号"`
	CreatedAt  []*gtime.Time `json:"createdAt"  dc:"发生时间"`
}

// LogListModel 资产流水
type LogListModel struct {
	entity.SysAssetLog
}

// StatInp 资产统计
type StatInp struct {
	UserId uint64 `json:"userId" dc:"用户ID，为空时统计全部"`
}

// StatModel 资产统计
type StatModel struct {
	TotalAvailable float64 `json:"totalAvailable" dc:"可用余额合计"`
	TotalFrozen    float64 `json:"totalFrozen"    dc:"冻结金额合计"`
	AccountCount   int     `json:"accountCount"   dc:"账户数"`
}
