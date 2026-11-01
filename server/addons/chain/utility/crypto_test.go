package utility

import (
	"context"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcfg"
)

func TestPrivateKeyEncryptDecrypt(t *testing.T) {
	adapter, err := gcfg.NewAdapterContent(`chain:
  privateKeySecret: "unit-test-private-key-secret"
`)
	if err != nil {
		t.Fatal(err)
	}
	g.Cfg().SetAdapter(adapter)

	ctx := context.Background()
	privateKey := "4c0883a69102937d6231471b5dbb6204fe5129617082792ae468d01a3f362318"

	cipherText, err := EncryptPrivateKey(ctx, privateKey)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}
	if !IsEncryptedPrivateKey(cipherText) {
		t.Fatalf("cipher text should has encrypted prefix: %s", cipherText)
	}
	if strings.Contains(cipherText, privateKey) {
		t.Fatal("cipher text should not contain plaintext private key")
	}

	plain, err := DecryptPrivateKey(ctx, cipherText)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}
	if plain != privateKey {
		t.Fatalf("decrypt mismatch, want: %s, got: %s", privateKey, plain)
	}

	// 已加密的私钥重复加密时应原样返回
	again, err := EncryptPrivateKey(ctx, cipherText)
	if err != nil {
		t.Fatalf("encrypt again failed: %v", err)
	}
	if again != cipherText {
		t.Fatal("encrypted private key should not be encrypted again")
	}

	// 明文私钥禁止解密
	if _, err = DecryptPrivateKey(ctx, privateKey); err == nil {
		t.Fatal("decrypt plaintext private key should return error")
	}
}
