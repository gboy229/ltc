// Package sysin
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sysin

import (
	"hotgo/internal/model/input/assetin"
)

// AssetAccountListInp 资产账户列表
type AssetAccountListInp struct {
	assetin.AccountListInp
	Uid           string `json:"uid"           dc:"用户UID"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

// AssetAccountListModel 资产账户
type AssetAccountListModel struct {
	assetin.AccountListModel
	Uid           string `json:"uid"           dc:"用户UID"`
	Nickname      string `json:"nickname"      dc:"用户昵称"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

// AssetLogListInp 资金流水列表
type AssetLogListInp struct {
	assetin.LogListInp
	Uid           string `json:"uid"           dc:"用户UID"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

// AssetLogListModel 资金流水
type AssetLogListModel struct {
	assetin.LogListModel
	Uid           string `json:"uid"           dc:"用户UID"`
	Nickname      string `json:"nickname"      dc:"用户昵称"`
	WalletAddress string `json:"walletAddress" dc:"钱包地址"`
}

// AssetStatInp 资产统计
type AssetStatInp struct {
	assetin.StatInp
}

// AssetStatModel 资产统计
type AssetStatModel struct {
	assetin.StatModel
}
