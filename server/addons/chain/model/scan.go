// Package model
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package model

// ScanTransferEvent 扫块推送到队列的标准化 Transfer 事件
type ScanTransferEvent struct {
	Chain         string `json:"chain"         dc:"链名"`
	Driver        string `json:"driver"        dc:"链驱动"`
	Contract      string `json:"contract"      dc:"合约地址"`
	Symbol        string `json:"symbol"        dc:"币种"`
	Decimals      int    `json:"decimals"      dc:"精度"`
	From          string `json:"from"          dc:"转出地址（用户钱包）"`
	To            string `json:"to"            dc:"转入地址"`
	Amount        string `json:"amount"        dc:"raw 金额"`
	TxHash        string `json:"txHash"        dc:"交易哈希"`
	LogIndex      int    `json:"logIndex"      dc:"日志序号"`
	BlockHeight   uint64 `json:"blockHeight"   dc:"块高"`
	BlockHash     string `json:"blockHash"     dc:"区块哈希"`
	Confirmations int    `json:"confirmations" dc:"确认数"`
}
