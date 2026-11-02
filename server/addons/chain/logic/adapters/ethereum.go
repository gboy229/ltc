package adapters

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/crypto/sha3"
	"hotgo/addons/chain/consts"
	"hotgo/addons/chain/utility"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"hotgo/internal/library/queue"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/goccy/go-json"
	"github.com/gogf/gf/v2/frame/g"
)

type EthereumAdapter struct {
	client  *ethclient.Client
	nodeURL string
	chain   string
}

// NewEthereumAdapter 创建一个新的以太坊适配器
func NewEthereumAdapter() (*EthereumAdapter, error) {
	return &EthereumAdapter{}, nil
}

func (ea *EthereumAdapter) Client(nodeURL string) error {
	client, err := ethclient.Dial(nodeURL) // 连接以太坊节点
	if err != nil {
		return fmt.Errorf("failed to connect to Ethereum node: %v", err)
	}
	ea.client = client
	ea.nodeURL = nodeURL
	return nil
}
func (ea *EthereumAdapter) SendTransaction(options TransactionOptions) (string, uint64, error) {

	if options.Amount <= 0 {
		return "", 0, fmt.Errorf("invalid amount: must be greater than zero")
	}

	if options.ToAddress == "" {
		return "", 0, fmt.Errorf("invalid address: 'ToAddress' cannot be empty")
	}

	if options.PrivateKey == "" {
		return "", 0, fmt.Errorf("invalid key: 'PrivateKey' cannot be empty")
	}

	var amount *big.Int
	var data []byte

	toAddress := common.HexToAddress(options.ToAddress)

	if options.TokenAddress == "" {
		amount = utility.ToWei(options.Amount, int32(ea.GetGasDecimal()))
		data = common.Hex2Bytes("0x")
	} else { //代币
		amount = big.NewInt(0)

		transferFnSignature := []byte("transfer(address,uint256)")
		hash := sha3.NewLegacyKeccak256()
		hash.Write(transferFnSignature)
		methodID := hash.Sum(nil)[:4]

		paddedAddress := common.LeftPadBytes(common.HexToAddress(options.ToAddress).Bytes(), 32)

		paddedAmount := common.LeftPadBytes(utility.ToWei(options.Amount, int32(options.Decimal)).Bytes(), 32)

		data = append(data, methodID...)
		data = append(data, paddedAddress...)
		data = append(data, paddedAmount...)

		toAddress = common.HexToAddress(options.TokenAddress)
	}

	gasPrice, err := ea.client.SuggestGasPrice(context.Background())

	if err != nil {
		return "", 0, fmt.Errorf("%v", err)
	}
	gasLimit, err := ea.client.EstimateGas(context.Background(), ethereum.CallMsg{})

	chainID, err := ea.client.NetworkID(context.Background())
	if err != nil {
		return "", 0, fmt.Errorf("%v", err)
	}

	privateKey, err := crypto.HexToECDSA(options.PrivateKey)
	if err != nil {
		return "", 0, fmt.Errorf("invalid private key: %v", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", 0, fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	var nonce uint64
	if options.Nonce == nil {
		nonce, err = ea.client.PendingNonceAt(context.Background(), fromAddress)
		if err != nil {
			return "", 0, fmt.Errorf("invalid nonce: %v", err)
		}
	} else {
		nonce = *options.Nonce
	}

	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      gasLimit,
		To:       &toAddress,
		Value:    amount,
		Data:     data,
	})

	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), privateKey)
	if err != nil {
		return "", 0, fmt.Errorf("%v", err)
	}

	err = ea.client.SendTransaction(context.Background(), signedTx)

	if err != nil {
		return "", 0, fmt.Errorf("%v", err)
	}
	return signedTx.Hash().String(), signedTx.Nonce(), nil
}

func (ea *EthereumAdapter) BlockTransaction(height uint64) (err error) {

	ctx := context.Background()

	blockNumber := new(big.Int).SetUint64(height)

	txList, err := ea.client.BlockByNumber(context.Background(), blockNumber)

	if err != nil && !errors.Is(err, ethereum.NotFound) {
		g.Log().Error(ctx, err)
		return err
	}

	if errors.Is(err, ethereum.NotFound) {
		return
	}

	for _, tx := range txList.Transactions() {

		dataStr := strings.ToLower(hex.EncodeToString(tx.Data()))

		if tx.To() == nil || (dataStr != "" && (len(dataStr) != 136 || dataStr[0:8] != "a9059cbb")) {
			continue
		}

		// 检查交易的ChainId是否为空或为0，如果是则跳过此交易
		txChainId := tx.ChainId()
		if txChainId == nil || txChainId.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		msg, err := core.TransactionToMessage(tx, types.LatestSignerForChainID(txChainId), nil)
		if err != nil {
			log.Printf("处理交易消息失败: %v, 跳过此交易", err)
			continue
		}

		var tokenAddress, toAddress, fromAddress string

		var value = new(big.Int)

		fromAddress = msg.From.Hex()
		if dataStr == "" {
			toAddress = tx.To().Hex()
			tokenAddress = ""
			value.SetString(tx.Value().String(), 10)
		} else {
			tokenAddress = tx.To().Hex()
			toAddress = dataStr[32:72]
			toAddress = fmt.Sprintf("0x%s", toAddress)
			value.SetString(dataStr[72:], 16)
		}

		exists, err := g.Redis().SIsMember(ctx, consts.AddressCacheListKey, strings.ToLower(toAddress))
		if err != nil {
			continue
		}

		if exists == 0 {
			continue
		}

		receipt, err := ea.client.TransactionReceipt(context.Background(), tx.Hash())
		if err != nil {
			continue
		}
		if receipt.Status == 1 {

			msg := NotifyMessage{
				Drive:        strings.ToLower("evm"),
				Chain:        ea.chain,
				Hash:         tx.Hash().Hex(),
				TokenAddress: tokenAddress,
				FromAddress:  fromAddress,
				ToAddress:    toAddress,
				Amount:       value.String(),
				Height:       txList.Number().String(),
			}

			pushData, err := json.Marshal(msg)
			if err != nil {
				log.Println(err)
				return err
			}

			// 使用系统队列推送消息
			err = queue.Push("scan_data", pushData)
			if err != nil {
				return err
			}

		}

	}

	return nil
}

func (ea *EthereumAdapter) Receipt(hash string) (bool, error) {
	ht := common.HexToHash(hash)

	_, isPending, err := ea.client.TransactionByHash(context.Background(), ht)
	if err != nil {
		return false, err
	}

	// 交易尚未上链，返回明确的 pending 错误，调用方据此跳过、等待下次确认
	if isPending {
		return false, fmt.Errorf("transaction is pending")
	}

	res, err := ea.client.TransactionReceipt(context.Background(), ht)
	if err != nil {
		return false, err
	}
	return res.Status == 1, nil
}

func (ea *EthereumAdapter) GetGasDecimal() int {
	return 18
}

func (ea *EthereumAdapter) GetBlockLastHeight() (string, error) {
	header, err := ea.client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return "", err
	}
	return header.Number.String(), nil
}

func (ea *EthereumAdapter) GetEstimateGas() (float64, error) {

	var amount float64

	gasPrice, err := ea.client.SuggestGasPrice(context.Background())

	if err != nil {
		return amount, err
	}

	gasLimit, err := ea.client.EstimateGas(context.Background(), ethereum.CallMsg{})

	if err != nil {
		return amount, err
	}
	gasAmount, _ := utility.ToDecimal(strconv.FormatUint(gasPrice.Uint64()*gasLimit, 10), 18).Round(8).Float64()

	gasAmount, _ = decimal.NewFromFloat(gasAmount * 1.5).Round(8).Float64()

	return gasAmount, nil
}

func (ea *EthereumAdapter) GetBalance(address string, decimal int) (float64, error) {
	var amount float64
	account := common.HexToAddress(address)
	balance, err := ea.client.BalanceAt(context.Background(), account, nil)
	if err != nil {
		return amount, fmt.Errorf("failed to get balance: %v", err)
	}

	amount, _ = utility.ToDecimal(balance, decimal).Float64()

	return amount, nil
}

type RPCRequest struct {
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      int           `json:"id"`
}
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type RPCResponse struct {
	Result string    `json:"result"`
	Error  *RPCError `json:"error,omitempty"`
	ID     int       `json:"id"`
}

func (ea *EthereumAdapter) GetTokenBalance(address, tokenAddress string, decimal int) (float64, error) {
	var amount float64

	if ea.nodeURL == "" {
		return amount, fmt.Errorf("node url is empty")
	}
	if !common.IsHexAddress(address) {
		return amount, fmt.Errorf("invalid wallet address: %s", address)
	}
	if !common.IsHexAddress(tokenAddress) {
		return amount, fmt.Errorf("invalid token address: %s", tokenAddress)
	}

	// balanceOf(address) 的 calldata：4 字节选择器 + 32 字节右对齐地址
	data := "0x70a08231" + strings.Repeat("0", 24) + strings.ToLower(strings.TrimPrefix(address, "0x"))

	// 构建请求
	request := RPCRequest{
		JSONRPC: "2.0",
		Method:  "eth_call",
		Params: []interface{}{
			map[string]interface{}{
				"to":   tokenAddress,
				"data": data,
			},
			"latest",
		},
	}

	// 将请求编码为 JSON
	reqBody, err := json.Marshal(request)
	if err != nil {
		return amount, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	// 发送 HTTP POST 请求
	resp, err := client.Post(ea.nodeURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return amount, err
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return amount, err
	}

	// 解析响应
	var rpcResponse RPCResponse
	err = json.Unmarshal(body, &rpcResponse)
	if err != nil {
		return amount, err
	}

	// 检查是否有错误
	if rpcResponse.Error != nil {
		return amount, fmt.Errorf("%s [Method]", rpcResponse.Error.Message)
	}
	if rpcResponse.Result == "" || rpcResponse.Result == "0x" {
		return amount, fmt.Errorf("empty balance result from token contract: %s", tokenAddress)
	}

	number, ok := new(big.Int).SetString(strings.TrimPrefix(strings.ToLower(rpcResponse.Result), "0x"), 16)
	if !ok {
		return amount, fmt.Errorf("invalid balance result: %s", rpcResponse.Result)
	}

	amount, _ = utility.ToDecimal(number.String(), decimal).Float64()

	return amount, nil
}

// GenerateAddress 生成新的以太坊地址
func (ea *EthereumAdapter) GenerateAddress() (address, private string, err error) {
	privateKey, err := crypto.GenerateKey() // 生成私钥
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", "", fmt.Errorf("cannot assert type: publicKey is not of type *ecdsa.PublicKe")
	}
	address = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	return address, hexutil.Encode(privateKeyBytes)[2:], nil

}

func (ea *EthereumAdapter) IsValidAddress(iaddress string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(iaddress)
}

func (ea *EthereumAdapter) Chain(chain string) {
	ea.chain = chain
}

var erc20TransferTopic = crypto.Keccak256Hash([]byte("Transfer(address,address,uint256)"))

// GetBlockNumber 获取链上最新块高
func (ea *EthereumAdapter) GetBlockNumber(ctx context.Context) (uint64, error) {
	if ea.client == nil {
		return 0, fmt.Errorf("client is not initialized")
	}
	return ea.client.BlockNumber(ctx)
}

// GetBlockHash 获取指定块高的区块哈希
func (ea *EthereumAdapter) GetBlockHash(ctx context.Context, height uint64) (string, error) {
	if ea.client == nil {
		return "", fmt.Errorf("client is not initialized")
	}
	header, err := ea.client.HeaderByNumber(ctx, new(big.Int).SetUint64(height))
	if err != nil {
		return "", fmt.Errorf("get block header %d: %v", height, err)
	}
	return header.Hash().Hex(), nil
}

// ScanTransferLogs 扫描指定块高区间内、指定合约的全部标准 ERC20 Transfer 日志
func (ea *EthereumAdapter) ScanTransferLogs(ctx context.Context, fromBlock, toBlock uint64, contracts []string) ([]*Transfer, error) {
	if ea.client == nil {
		return nil, fmt.Errorf("client is not initialized")
	}

	addresses := make([]common.Address, 0, len(contracts))
	for _, contract := range contracts {
		if !common.IsHexAddress(contract) {
			continue
		}
		addresses = append(addresses, common.HexToAddress(contract))
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("no valid contract address: %v", contracts)
	}

	return ea.scanTransferLogs(ctx, fromBlock, toBlock, addresses)
}

// scanTransferLogs 真正执行 eth_getLogs；遇到节点范围/结果条数限制时二分区间重试
func (ea *EthereumAdapter) scanTransferLogs(ctx context.Context, fromBlock, toBlock uint64, addresses []common.Address) ([]*Transfer, error) {
	query := ethereum.FilterQuery{
		FromBlock: new(big.Int).SetUint64(fromBlock),
		ToBlock:   new(big.Int).SetUint64(toBlock),
		Addresses: addresses,
		Topics: [][]common.Hash{
			{erc20TransferTopic},
		},
	}

	logs, err := ea.client.FilterLogs(ctx, query)
	if err != nil {
		if isRangeLimitError(err) {
			if fromBlock == toBlock {
				return nil, fmt.Errorf("filter logs %d: %v（当前 RPC 不支持 eth_getLogs 或已限流，请更换支持日志查询的 RPC）", fromBlock, err)
			}
			mid := fromBlock + (toBlock-fromBlock)/2
			left, leftErr := ea.scanTransferLogs(ctx, fromBlock, mid, addresses)
			if leftErr != nil {
				return nil, leftErr
			}
			right, rightErr := ea.scanTransferLogs(ctx, mid+1, toBlock, addresses)
			if rightErr != nil {
				return nil, rightErr
			}
			return append(left, right...), nil
		}
		return nil, fmt.Errorf("filter logs %d-%d: %v", fromBlock, toBlock, err)
	}

	transfers := make([]*Transfer, 0, len(logs))
	for _, item := range logs {
		if item.Removed {
			continue
		}
		// 标准 ERC20 Transfer：3 个 topic（事件签名/from/to）+ 32 字节 value。
		// ERC721 的 Transfer 有 4 个 topic（tokenId 也 indexed），这里跳过。
		if len(item.Topics) != 3 || len(item.Data) != 32 {
			continue
		}

		to := common.HexToAddress(item.Topics[2].Hex())
		from := common.HexToAddress(item.Topics[1].Hex())
		amount := new(big.Int).SetBytes(item.Data).String()

		transfers = append(transfers, &Transfer{
			Contract:    item.Address.Hex(),
			From:        from.Hex(),
			To:          to.Hex(),
			Amount:      amount,
			TxHash:      item.TxHash.Hex(),
			LogIndex:    int(item.Index),
			BlockHeight: item.BlockNumber,
			BlockHash:   item.BlockHash.Hex(),
		})
	}
	return transfers, nil
}

// isRangeLimitError 判断是否为 RPC 的范围/结果条数/超时限制，可尝试二分区间
func isRangeLimitError(err error) bool {
	msg := strings.ToLower(err.Error())
	for _, key := range []string{"limit", "too many", "more than", "exceed", "too large", "query timeout", "timeout"} {
		if strings.Contains(msg, key) {
			return true
		}
	}
	return false
}
