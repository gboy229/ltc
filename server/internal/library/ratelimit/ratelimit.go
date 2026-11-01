// Package ratelimit 基于 Redis 的简单操作限流。
package ratelimit

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

// Allow 在 ttl 内只允许同一个 key 成功一次。
func Allow(ctx context.Context, key string, ttl time.Duration) (bool, error) {
	value, err := g.Redis().Do(ctx, "SET", key, 1, "NX", "EX", int(ttl.Seconds()))
	if err != nil {
		return false, gerror.Wrap(err, "限流校验失败")
	}
	return !value.IsNil() && value.String() == "OK", nil
}

// AllowCount 在时间窗口内限制同一个 key 的最大请求次数。
func AllowCount(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	if limit <= 0 {
		return false, gerror.New("限流阈值必须大于0")
	}
	const script = `local n = redis.call("INCR", KEYS[1])
if n == 1 then redis.call("PEXPIRE", KEYS[1], ARGV[2]) end
if n > tonumber(ARGV[1]) then return 0 end
return 1`
	value, err := g.Redis().Do(ctx, "EVAL", script, 1, key, limit, int(window.Milliseconds()))
	if err != nil {
		return false, gerror.Wrap(err, "限流校验失败")
	}
	return value.Int() == 1, nil
}
