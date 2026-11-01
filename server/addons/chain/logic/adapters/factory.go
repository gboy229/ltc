package adapters

import (
	"fmt"
	"strings"
)

// NewChainAdapter 是一个工厂函数，用于根据链类型创建相应的适配器
func NewChainAdapter(chainType string) (ChainAdapter, error) {
	switch strings.ToLower(chainType) {
	case "evm":
		return NewEthereumAdapter()
	default:
		return nil, fmt.Errorf("unsupported chain type: %s", chainType)
	}
}
