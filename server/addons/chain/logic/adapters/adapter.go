package adapters

import "context"

type ChainAdapter interface {
	Client(nodeUrl string) error
	SendTransaction(options TransactionOptions) (string, uint64, error)
	BlockTransaction(height uint64) error
	GetBlockLastHeight() (string, error)
	GetBalance(address string, decimal int) (float64, error)
	GetTokenBalance(address, tokenAddress string, decimal int) (float64, error)
	GetGasDecimal() int
	GetEstimateGas() (float64, error)
	GenerateAddress() (string, string, error)
	Receipt(hash string) (bool, error)
	IsValidAddress(address string) bool
	Chain(name string)
	// GetBlockNumber 获取链上最新块高
	GetBlockNumber(ctx context.Context) (uint64, error)
	// GetBlockHash 获取指定块高的区块哈希
	GetBlockHash(ctx context.Context, height uint64) (string, error)
	// ScanTransferLogs 扫描指定块高区间内、指定合约的全部标准 Transfer 日志
	ScanTransferLogs(ctx context.Context, fromBlock, toBlock uint64, contracts []string) ([]*Transfer, error)
}

// Transfer 标准 Transfer 事件（ERC20）
type Transfer struct {
	Contract    string `json:"contract"`    // 合约地址
	From        string `json:"from"`        // 转出地址
	To          string `json:"to"`          // 转入地址
	Amount      string `json:"amount"`      // raw 金额
	TxHash      string `json:"txHash"`      // 交易哈希
	LogIndex    int    `json:"logIndex"`    // 日志序号
	BlockHeight uint64 `json:"blockHeight"` // 块高
	BlockHash   string `json:"blockHash"`   // 区块哈希
}

type TransactionOptions struct {
	PrivateKey   string
	FromAddress  string
	ToAddress    string
	TokenAddress string
	Amount       float64
	IsAccelerate bool
	Decimal      int8
	Nonce        *uint64
}

type NotifyMessage struct {
	Drive        string `json:"drive"`
	Chain        string `json:"chain"`
	Hash         string `json:"hash"`
	TokenAddress string `json:"token_address"`
	FromAddress  string `json:"from_address"`
	ToAddress    string `json:"to_address"`
	Height       string `json:"height"`
	Amount       string `json:"amount"`
}
