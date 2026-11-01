// Package utility
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2024 HotGo CLI
// @Author  Ms <133814250@qq.com>
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package utility

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/decred/dcrd/dcrec/secp256k1/v4/ecdsa"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/util/grand"
	"golang.org/x/crypto/sha3"
)

// IsValidWalletAddress 验证 EVM 钱包地址是否合法（0x + 40 位十六进制，大小写不敏感）
func IsValidWalletAddress(address string) bool {
	if len(address) != 42 || !strings.HasPrefix(address, "0x") {
		return false
	}
	_, err := hex.DecodeString(address[2:])
	return err == nil
}

// VerifyWalletSignature 验证 EVM 钱包 personal_sign 签名，返回签名地址与给定地址是否一致
func VerifyWalletSignature(address, message, signatureHex string) (bool, error) {
	sig, err := hex.DecodeString(strings.TrimPrefix(signatureHex, "0x"))
	if err != nil {
		return false, gerror.New("签名格式错误")
	}
	if len(sig) != 65 {
		return false, gerror.New("签名长度错误")
	}

	// Ethereum 签名格式：[R(32)][S(32)][v]，v 为 27/28
	v := sig[64]
	if v != 27 && v != 28 {
		return false, gerror.New("签名恢复标识非法")
	}

	hash := personalSignHash([]byte(message))

	// 转为 dcrec RecoverCompact 期望的格式：[recoveryCode(27/28)][R(32)][S(32)]
	recoverySig := make([]byte, 65)
	recoverySig[0] = v
	copy(recoverySig[1:33], sig[0:32])
	copy(recoverySig[33:65], sig[32:64])

	pubKey, _, err := ecdsa.RecoverCompact(recoverySig, hash)
	if err != nil {
		return false, gerror.New("签名验证失败")
	}

	addressBytes := keccak256(pubKey.SerializeUncompressed()[1:])[12:]
	recovered := "0x" + hex.EncodeToString(addressBytes)

	return strings.EqualFold(recovered, address), nil
}

// personalSignHash 计算 personal_sign 的消息哈希（EIP-191）
func personalSignHash(message []byte) []byte {
	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(message))
	return keccak256(append([]byte(prefix), message...))
}

// keccak256 计算 Keccak-256 哈希
func keccak256(data []byte) []byte {
	h := sha3.NewLegacyKeccak256()
	_, _ = h.Write(data)
	return h.Sum(nil)
}

// GenUid 生成对外唯一用户号
func GenUid() string {
	return grand.S(32)
}
