package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEthereumAdapterScanTransferLogs(t *testing.T) {
	const (
		contract  = "0x55d398326f99059ff775485246999027b3197955"
		from      = "0xeb2c26a1b36f16a8506d334119beb6e147076517"
		to        = "0x5ec081cf77c0796545db0b89101b2ec2efc21bd6"
		otherTo   = "0x4d9ddacfa40fb42fa65ebb6293f758e52dc6c3cf"
		txHash    = "0xc272ed78acc63f51717270e82e692be00f6aa3f218cc349ae7960186cb97ae94"
		blockHash = "0x1111111111111111111111111111111111111111111111111111111111111111"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			ID     int             `json:"id"`
			Method string          `json:"method"`
			Params json.RawMessage `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}

		var result interface{}
		switch req.Method {
		case "eth_blockNumber":
			result = "0x10"
		case "eth_getBlockByNumber":
			result = map[string]interface{}{
				"parentHash":       "0x0000000000000000000000000000000000000000000000000000000000000001",
				"sha3Uncles":       "0x1dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347",
				"miner":            "0x0000000000000000000000000000000000000000",
				"stateRoot":        "0x0000000000000000000000000000000000000000000000000000000000000002",
				"transactionsRoot": "0x0000000000000000000000000000000000000000000000000000000000000003",
				"receiptsRoot":     "0x0000000000000000000000000000000000000000000000000000000000000004",
				"logsBloom":        "0x" + strings.Repeat("0", 512),
				"difficulty":       "0x1",
				"number":           "0xf",
				"gasLimit":         "0x1",
				"gasUsed":          "0x1",
				"timestamp":        "0x1",
				"extraData":        "0x",
				"mixHash":          "0x0000000000000000000000000000000000000000000000000000000000000005",
				"nonce":            "0x0000000000000000",
				"baseFeePerGas":    "0x1",
			}
		case "eth_getLogs":
			result = []map[string]interface{}{
				{
					"address": contract,
					"topics": []string{
						"0xddf252ad1be25ef4a7c1e5b8f5f0d7c6a2b1e8c9d0f3a4b5c6d7e8f9a0b1c2d3",
						"0x000000000000000000000000eb2c26a1b36f16a8506d334119beb6e147076517",
						"0x0000000000000000000000005ec081cf77c0796545db0b89101b2ec2efc21bd6",
					},
					"data":             "0x0000000000000000000000000000000000000000000000000000000000000064",
					"blockNumber":      "0x10",
					"transactionHash":  txHash,
					"transactionIndex": "0x1",
					"blockHash":        blockHash,
					"logIndex":         "0x2",
					"removed":          false,
				},
				{
					"address": contract,
					"topics": []string{
						"0xddf252ad1be25ef4a7c1e5b8f5f0d7c6a2b1e8c9d0f3a4b5c6d7e8f9a0b1c2d3",
						"0x000000000000000000000000eb2c26a1b36f16a8506d334119beb6e147076517",
						"0x0000000000000000000000004d9ddacfa40fb42fa65ebb6293f758e52dc6c3cf",
					},
					"data":             "0x00000000000000000000000000000000000000000000000000000000000000c8",
					"blockNumber":      "0x10",
					"transactionHash":  txHash,
					"transactionIndex": "0x1",
					"blockHash":        blockHash,
					"logIndex":         "0x3",
					"removed":          false,
				},
			}
		default:
			result = nil
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      req.ID,
			"result":  result,
		})
	}))
	defer server.Close()

	adapter, err := NewChainAdapter("evm")
	if err != nil {
		t.Fatal(err)
	}
	if err = adapter.Client(server.URL); err != nil {
		t.Fatal(err)
	}

	latest, err := adapter.GetBlockNumber(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if latest != 16 {
		t.Fatalf("latest block mismatch, want 16, got %d", latest)
	}

	gotBlockHash, err := adapter.GetBlockHash(context.Background(), 15)
	if err != nil {
		t.Fatal(err)
	}
	if gotBlockHash == "" {
		t.Fatal("block hash is empty")
	}

	transfers, err := adapter.ScanTransferLogs(context.Background(), 1, 16, []string{contract})
	if err != nil {
		t.Fatal(err)
	}
	if len(transfers) != 2 {
		t.Fatalf("transfers length mismatch, want 2, got %d", len(transfers))
	}

	var got, other *Transfer
	for _, item := range transfers {
		switch item.LogIndex {
		case 2:
			got = item
		case 3:
			other = item
		}
	}
	if got == nil || other == nil {
		t.Fatalf("missing expected transfers: %+v", transfers)
	}
	if !strings.EqualFold(other.To, otherTo) {
		t.Fatalf("other to mismatch: %s", other.To)
	}
	if other.Amount != "200" {
		t.Fatalf("other amount mismatch: %s", other.Amount)
	}
	if !strings.EqualFold(got.Contract, contract) {
		t.Fatalf("contract mismatch: %s", got.Contract)
	}
	if !strings.EqualFold(got.From, from) {
		t.Fatalf("from mismatch: %s", got.From)
	}
	if !strings.EqualFold(got.To, to) {
		t.Fatalf("to mismatch: %s", got.To)
	}
	if got.Amount != "100" {
		t.Fatalf("amount mismatch: %s", got.Amount)
	}
	if got.TxHash != txHash {
		t.Fatalf("tx hash mismatch: %s", got.TxHash)
	}
	if got.LogIndex != 2 {
		t.Fatalf("log index mismatch: %d", got.LogIndex)
	}
	if got.BlockHeight != 16 {
		t.Fatalf("block height mismatch: %d", got.BlockHeight)
	}
}
