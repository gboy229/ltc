package adapters

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// TestEthereumAdapterScanTransferLogsSplitsRange 验证 RPC 返回 limit exceeded 时会二分区间继续扫
func TestEthereumAdapterScanTransferLogsSplitsRange(t *testing.T) {
	const (
		contract = "0x55d398326f99059ff775485246999027b3197955"
		fromAddr = "0xeb2c26a1b36f16a8506d334119beb6e147076517"
		toAddr   = "0x5ec081cf77c0796545db0b89101b2ec2efc21bd6"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     int `json:"id"`
			Params []struct {
				FromBlock string `json:"fromBlock"`
				ToBlock   string `json:"toBlock"`
			} `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}

		var resp map[string]interface{}
		if len(req.Params) == 0 {
			resp = map[string]interface{}{"jsonrpc": "2.0", "id": req.ID, "result": []interface{}{}}
		} else {
			fromBlock := parseHexUint(req.Params[0].FromBlock)
			toBlock := parseHexUint(req.Params[0].ToBlock)

			if toBlock-fromBlock > 10 {
				resp = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req.ID,
					"error": map[string]interface{}{
						"code":    -32005,
						"message": "limit exceeded",
					},
				}
			} else {
				resp = map[string]interface{}{
					"jsonrpc": "2.0",
					"id":      req.ID,
					"result": []map[string]interface{}{
						{
							"address": contract,
							"topics": []string{
								"0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef",
								"0x000000000000000000000000" + strings.TrimPrefix(fromAddr, "0x"),
								"0x000000000000000000000000" + strings.TrimPrefix(toAddr, "0x"),
							},
							"data":             "0x0000000000000000000000000000000000000000000000000000000000000064",
							"blockNumber":      fmt.Sprintf("0x%x", fromBlock),
							"transactionHash":  "0x" + strings.Repeat("11", 32),
							"transactionIndex": "0x0",
							"blockHash":        "0x" + strings.Repeat("22", 32),
							"logIndex":         "0x0",
							"removed":          false,
						},
					},
				}
			}
		}

		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	adapter, err := NewChainAdapter("evm")
	if err != nil {
		t.Fatal(err)
	}
	if err = adapter.Client(server.URL); err != nil {
		t.Fatal(err)
	}

	transfers, err := adapter.ScanTransferLogs(context.Background(), 1, 100, []string{contract})
	if err != nil {
		t.Fatal(err)
	}
	if len(transfers) == 0 {
		t.Fatal("expected transfers after splitting range, got 0")
	}
	for _, transfer := range transfers {
		if transfer.BlockHeight < 1 || transfer.BlockHeight > 100 {
			t.Fatalf("unexpected block height: %d", transfer.BlockHeight)
		}
	}
}

func parseHexUint(s string) uint64 {
	v, _ := strconv.ParseUint(strings.TrimPrefix(s, "0x"), 16, 64)
	return v
}
