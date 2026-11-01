// Package walletsign 钱包 nonce 与签名校验。
package walletsign

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/grand"

	usersconsts "hotgo/addons/users/consts"
	"hotgo/addons/users/utility"
)

// Generate 生成钱包签名挑战值。
func Generate(ctx context.Context, address string) (string, error) {
	address = normalize(address)
	nonce := grand.S(32)
	if _, err := g.Redis().Do(ctx, "SET", usersconsts.CacheWalletNonce+address, nonce, "EX", 300); err != nil {
		return "", gerror.Wrap(err, "生成签名挑战值失败")
	}
	return nonce, nil
}

// VerifyAndConsume 校验签名并原子消费 nonce。
func VerifyAndConsume(ctx context.Context, address, nonce, signature string) error {
	address = normalize(address)
	key := usersconsts.CacheWalletNonce + address
	nonceVal, err := g.Redis().Do(ctx, "GET", key)
	if err != nil {
		return gerror.Wrap(err, "获取签名挑战值失败")
	}
	if nonceVal.IsNil() || nonceVal.IsEmpty() {
		return gerror.New("签名挑战值已失效，请重新获取")
	}
	if nonceVal.String() != nonce {
		return gerror.New("签名挑战值不匹配")
	}

	ok, err := utility.VerifyWalletSignature(address, nonce, signature)
	if err != nil {
		return err
	}
	if !ok {
		return gerror.New("签名验证失败")
	}

	const script = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]) else return 0 end`
	value, err := g.Redis().Do(ctx, "EVAL", script, 1, key, nonce)
	if err != nil {
		return gerror.Wrap(err, "消费签名挑战值失败")
	}
	if value.Int() != 1 {
		return gerror.New("签名挑战值已失效，请重新获取")
	}
	return nil
}

func normalize(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}
