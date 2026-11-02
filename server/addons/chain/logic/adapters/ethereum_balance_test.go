package adapters

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestGetTokenBalance 校验 balanceOf 的 calldata 编码与返回值解析
func TestGetTokenBalance(t *testing.T) {
	const (
		wallet   = "0x295379a4f4fc5f9ddb4fb67a85b53bca083bed5e"
		contract = "0x55d398326f99059ff775485246999027b3197955"
		expected = "0x70a08231000000000000000000000000295379a4f4fc5f9ddb4fb67a85b53bca083bed5e"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Method string        `json:"method"`
			Params []interface{} `json:"params"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("decode request: %v", err)
			return
		}
		if req.Method != "eth_call" {
			t.Errorf("unexpected method: %s", req.Method)
		}
		if len(req.Params) != 2 {
			t.Errorf("unexpected params length: %d", len(req.Params))
			return
		}

		call, ok := req.Params[0].(map[string]interface{})
		if !ok {
			t.Errorf("unexpected call params: %#v", req.Params[0])
			return
		}
		if call["to"] != contract {
			t.Errorf("unexpected to: %v", call["to"])
		}
		if call["data"] != expected {
			t.Errorf("unexpected calldata:\n got: %v\nwant: %v", call["data"], expected)
		}

		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"jsonrpc": "2.0",
			"id":      1,
			"result":  "0x0000000000000000000000000000000000000000000000004e65487ad4d39000",
		})
	}))
	defer server.Close()

	adapter := &EthereumAdapter{nodeURL: server.URL}
	amount, err := adapter.GetTokenBalance(wallet, contract, 18)
	if err != nil {
		t.Fatalf("GetTokenBalance error: %v", err)
	}
	if amount != 5.649001 {
		t.Fatalf("unexpected amount: %v", amount)
	}
}
