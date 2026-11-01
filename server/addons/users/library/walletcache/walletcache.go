// Package walletcache 钱包地址与用户ID映射。
// 注册/登录时写入 Redis，扫块时通过付款地址反查用户ID。
package walletcache

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// HashKey 钱包地址 -> 用户ID
const HashKey = "users:wallet:uid"

func normalize(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

// Set 写入钱包地址与用户ID映射。
func Set(ctx context.Context, address string, userId int64) error {
	address = normalize(address)
	if address == "" {
		return gerror.New("钱包地址不能为空")
	}
	if userId == 0 {
		return gerror.New("用户ID不能为空")
	}

	_, err := g.Redis().Do(ctx, "HSET", HashKey, address, userId)
	if err != nil {
		return gerror.Wrap(err, "写入钱包用户缓存失败")
	}
	return nil
}

// Get 根据钱包地址读取用户ID，不存在返回 0。
func Get(ctx context.Context, address string) (int64, error) {
	address = normalize(address)
	if address == "" {
		return 0, nil
	}

	value, err := g.Redis().Do(ctx, "HGET", HashKey, address)
	if err != nil {
		return 0, gerror.Wrap(err, "读取钱包用户缓存失败")
	}
	if value == nil || value.IsNil() || value.IsEmpty() {
		return 0, nil
	}
	return value.Int64(), nil
}

// Delete 删除钱包地址映射。
func Delete(ctx context.Context, address string) error {
	address = normalize(address)
	if address == "" {
		return nil
	}
	_, err := g.Redis().Do(ctx, "HDEL", HashKey, address)
	return err
}
