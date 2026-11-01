package model

// ChainPaymentEvent 链上支付确认事件
// chain 插件校验并落库 chain_recharge 后推送，业务插件幂等消费。
type ChainPaymentEvent struct {
	RechargeId  uint64 `json:"rechargeId"  dc:"充值记录ID"`
	Uid         int64  `json:"uid"         dc:"用户ID，0表示无法解析"`
	Chain       string `json:"chain"       dc:"链名"`
	Symbol      string `json:"symbol"      dc:"币种"`
	Contract    string `json:"contract"    dc:"合约地址"`
	From        string `json:"from"        dc:"付款地址"`
	To          string `json:"to"          dc:"收款地址"`
	Amount      string `json:"amount"      dc:"金额，十进制字符串"`
	TxHash      string `json:"txHash"      dc:"交易哈希"`
	LogIndex    int    `json:"logIndex"    dc:"日志序号"`
	BlockHeight int64  `json:"blockHeight" dc:"区块高度"`
}
