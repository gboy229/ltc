//go:build ignore

// 临时测试用户生成命令，测试完成后可直接删除本文件。
//
// 在 server 目录执行：
//
//	go run addons/users/testuser.go -n 10
//	go run addons/users/testuser.go -n 5 -parent 85285352
//
// 每个用户都会生成随机 EVM 钱包，调用真实登录接口：
//   - GET  /api/users/auth/nonce
//   - POST /api/users/auth/login
//
// 因此会完整走原有注册、钱包签名校验、推荐关系与 users_tree 写入逻辑。
package main

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/go-sql-driver/mysql"
)

var (
	baseURL string
	client  = &http.Client{Timeout: 15 * time.Second}
	db      *sql.DB
)

func main() {
	count := flag.Int("n", 10, "生成用户数量")
	parentUID := flag.String("parent", "", "固定上级UID，不传则随机绑定已有用户")
	dsn := flag.String("dsn", "root:root@tcp(127.0.0.1:3306)/zztest?charset=utf8mb4&parseTime=true&loc=Local", "MySQL DSN，仅用于查询可绑定上级")
	base := flag.String("url", "http://127.0.0.1:8000", "服务地址")
	flag.Parse()

	if *count <= 0 {
		fmt.Println("生成数量必须大于0")
		os.Exit(1)
	}
	baseURL = strings.TrimRight(*base, "/")

	var err error
	db, err = sql.Open("mysql", *dsn)
	if err != nil {
		fmt.Println("连接数据库失败:", err)
		os.Exit(1)
	}
	defer db.Close()
	if err = db.Ping(); err != nil {
		fmt.Println("数据库连接不可用:", err)
		os.Exit(1)
	}

	parentUIDs, err := loadParentUIDs(context.Background())
	if err != nil {
		fmt.Println("读取可绑定上级失败:", err)
		os.Exit(1)
	}

	for i := 0; i < *count; i++ {
		referrer := *parentUID
		if referrer == "" {
			if len(parentUIDs) == 0 {
				fmt.Println("没有可绑定的已有用户，请先手动创建或在空库时不传上级")
				os.Exit(1)
			}
			referrer = parentUIDs[randInt(len(parentUIDs))]
		}

		address, signature, nonce, err := walletSignature(context.Background())
		if err != nil {
			fmt.Println("生成钱包签名失败:", err)
			os.Exit(1)
		}

		token, err := login(context.Background(), address, nonce, signature, referrer)
		if err != nil {
			fmt.Println("登录失败:", err)
			os.Exit(1)
		}

		uid, id, err := findUserByAddress(context.Background(), address)
		if err != nil {
			fmt.Println("查询新用户失败:", err)
			os.Exit(1)
		}
		if uid != "" {
			parentUIDs = append(parentUIDs, uid)
		}
		fmt.Printf("created id=%d uid=%s parent=%s address=%s token=%s...\n", id, uid, referrer, address, token[:12])
	}
}

func loadParentUIDs(ctx context.Context) ([]string, error) {
	rows, err := db.QueryContext(ctx, "SELECT uid FROM hg_addon_users_user WHERE deleted_at IS NULL ORDER BY id ASC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	uids := make([]string, 0)
	for rows.Next() {
		var uid string
		if err = rows.Scan(&uid); err != nil {
			return nil, err
		}
		uids = append(uids, uid)
	}
	return uids, rows.Err()
}

func findUserByAddress(ctx context.Context, address string) (uid string, id int64, err error) {
	err = db.QueryRowContext(ctx, `
		SELECT u.id, u.uid
		FROM hg_addon_users_user u
		JOIN hg_addon_users_user_account a ON a.user_id = u.id
		WHERE a.type = 'wallet' AND a.account = ?
		LIMIT 1`, address).Scan(&id, &uid)
	if err == sql.ErrNoRows {
		return "", 0, nil
	}
	return uid, id, err
}

func walletSignature(ctx context.Context) (address, signature, nonce string, err error) {
	key, err := crypto.GenerateKey()
	if err != nil {
		return "", "", "", err
	}
	address = strings.ToLower(crypto.PubkeyToAddress(key.PublicKey).Hex())

	nonce, err = getNonce(ctx, address)
	if err != nil {
		return "", "", "", err
	}

	hash := accounts.TextHash([]byte(nonce))
	sig, err := crypto.Sign(hash, key)
	if err != nil {
		return "", "", "", err
	}
	if sig[64] < 27 {
		sig[64] += 27
	}
	return address, "0x" + hex.EncodeToString(sig), nonce, nil
}

func getNonce(ctx context.Context, address string) (string, error) {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Nonce string `json:"nonce"`
		} `json:"data"`
	}
	endpoint := fmt.Sprintf("%s/api/users/auth/nonce?address=%s", baseURL, url.QueryEscape(address))
	if err := requestJSON(ctx, http.MethodGet, endpoint, nil, &res); err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("%s", res.Message)
	}
	return res.Data.Nonce, nil
}

func login(ctx context.Context, address, nonce, signature, referrer string) (string, error) {
	var res struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
		Data    struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	body := map[string]string{
		"type":      "wallet",
		"account":   address,
		"nonce":     nonce,
		"signature": signature,
		"uid":       referrer,
	}
	if err := requestJSON(ctx, http.MethodPost, baseURL+"/api/users/auth/login", body, &res); err != nil {
		return "", err
	}
	if res.Code != 0 {
		return "", fmt.Errorf("%s", res.Message)
	}
	return res.Data.Token, nil
}

func requestJSON(ctx context.Context, method, endpoint string, body any, out any) error {
	var reader io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewReader(data)
	}
	req, err := http.NewRequestWithContext(ctx, method, endpoint, reader)
	if err != nil {
		return err
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(data))
	}
	return json.Unmarshal(data, out)
}

func randInt(max int) int {
	if max <= 0 {
		return 0
	}
	n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
	if err != nil {
		return 0
	}
	return int(n.Int64())
}
