package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]
local rate = tonumber(ARGV[1])
local burst = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local ttl = tonumber(ARGV[4])

local values = redis.call("HMGET", key, "tokens", "updated")
local tokens = tonumber(values[1])
local updated = tonumber(values[2])

if tokens == nil then
  tokens = burst
  updated = now
end

tokens = math.min(burst, tokens + ((now - updated) / 1000) * rate)
local allowed = 0
if tokens >= 1 then
  tokens = tokens - 1
  allowed = 1
end

redis.call("HSET", key, "tokens", tokens, "updated", now)
redis.call("PEXPIRE", key, ttl)
return allowed
`)

func Allow(ctx context.Context, key string, rate float64, burst int) (bool, error) {
	if rate <= 0 || burst <= 0 {
		return true, nil
	}
	ttl := time.Duration(float64(time.Second) * float64(burst) / rate)
	if ttl < time.Second {
		ttl = time.Second
	}
	result, err := tokenBucketScript.Run(
		ctx,
		client,
		[]string{key},
		rate,
		burst,
		time.Now().UnixMilli(),
		ttl.Milliseconds()*2,
	).Int()
	return result == 1, err
}
