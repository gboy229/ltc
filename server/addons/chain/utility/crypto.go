// Package utility
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package utility

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"os"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

const (
	// privateKeySecretConfig 私钥加密密钥的配置项
	privateKeySecretConfig = "chain.privateKeySecret"
	// privateKeySecretEnv 私钥加密密钥的环境变量，优先级高于配置项
	privateKeySecretEnv = "CHAIN_PRIVATE_KEY_SECRET"
	// privateKeyCipherPrefix 已加密私钥的标识前缀
	privateKeyCipherPrefix = "enc:v1:"
)

// privateKeySecret 获取私钥加密密钥，优先读取环境变量，其次读取配置项
func privateKeySecret(ctx context.Context) string {
	if v := strings.TrimSpace(os.Getenv(privateKeySecretEnv)); v != "" {
		return v
	}
	return g.Cfg().MustGet(ctx, privateKeySecretConfig, "").String()
}

// IsEncryptedPrivateKey 判断私钥是否已加密
func IsEncryptedPrivateKey(privateKey string) bool {
	return strings.HasPrefix(privateKey, privateKeyCipherPrefix)
}

// EncryptPrivateKey 加密钱包私钥，使用 AES-256-GCM
// 加密结果格式：enc:v1:{base64(nonce+ciphertext)}
func EncryptPrivateKey(ctx context.Context, privateKey string) (string, error) {
	if privateKey == "" {
		return "", gerror.New("钱包私钥不能为空")
	}

	// 已加密的数据直接返回，避免重复加密
	if IsEncryptedPrivateKey(privateKey) {
		return privateKey, nil
	}

	secret := privateKeySecret(ctx)
	if secret == "" {
		return "", gerror.Newf("未配置私钥加密密钥，请设置环境变量 %s 或配置项 %s", privateKeySecretEnv, privateKeySecretConfig)
	}

	gcm, err := newPrivateKeyGCM(secret)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", gerror.Wrap(err, "生成私钥加密随机数失败")
	}

	sealed := gcm.Seal(nonce, nonce, []byte(privateKey), nil)
	return privateKeyCipherPrefix + base64.StdEncoding.EncodeToString(sealed), nil
}

// DecryptPrivateKey 解密钱包私钥
func DecryptPrivateKey(ctx context.Context, cipherText string) (string, error) {
	if cipherText == "" {
		return "", gerror.New("钱包私钥不能为空")
	}
	if !IsEncryptedPrivateKey(cipherText) {
		return "", gerror.New("钱包私钥未加密，拒绝解密")
	}

	secret := privateKeySecret(ctx)
	if secret == "" {
		return "", gerror.Newf("未配置私钥加密密钥，请设置环境变量 %s 或配置项 %s", privateKeySecretEnv, privateKeySecretConfig)
	}

	gcm, err := newPrivateKeyGCM(secret)
	if err != nil {
		return "", err
	}

	raw, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(cipherText, privateKeyCipherPrefix))
	if err != nil {
		return "", gerror.New("钱包私钥密文格式错误")
	}

	nonceSize := gcm.NonceSize()
	if len(raw) < nonceSize {
		return "", gerror.New("钱包私钥密文格式错误")
	}

	plain, err := gcm.Open(nil, raw[:nonceSize], raw[nonceSize:], nil)
	if err != nil {
		return "", gerror.New("钱包私钥解密失败，请检查加密密钥是否正确")
	}
	return string(plain), nil
}

// newPrivateKeyGCM 创建 AES-256-GCM 加密器
func newPrivateKeyGCM(secret string) (cipher.AEAD, error) {
	// 使用 sha256 将任意长度的密钥整理为 32 字节
	key := sha256.Sum256([]byte(secret))
	block, err := aes.NewCipher(key[:])
	if err != nil {
		return nil, gerror.Wrap(err, "初始化私钥加密算法失败")
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, gerror.Wrap(err, "初始化私钥加密模式失败")
	}
	return gcm, nil
}
