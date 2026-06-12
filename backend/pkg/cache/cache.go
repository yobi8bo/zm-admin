package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"zhanxu-admin/backend/config"
)

var client *redis.Client
var keyPrefix = "zhanxu-admin:"
var defaultTTL = 15 * time.Minute
var captchaTTL = 5 * time.Minute
var rotateValueScript = redis.NewScript(`
if redis.call("GET", KEYS[1]) ~= ARGV[1] then
  return 0
end
redis.call("SET", KEYS[1], ARGV[2], "PX", ARGV[3])
return 1
`)

func Init(cfg config.RedisConfig) error {
	if cfg.KeyPrefix != "" {
		keyPrefix = cfg.KeyPrefix
	}
	if cfg.CacheTTL > 0 {
		defaultTTL = time.Duration(cfg.CacheTTL) * time.Second
	}
	if cfg.CaptchaExpire > 0 {
		captchaTTL = time.Duration(cfg.CaptchaExpire) * time.Second
	}
	client = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  durationOrDefault(cfg.DialTimeout, 5*time.Second),
		ReadTimeout:  durationOrDefault(cfg.ReadTimeout, 3*time.Second),
		WriteTimeout: durationOrDefault(cfg.WriteTimeout, 3*time.Second),
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func Get() *redis.Client { return client }

func DefaultTTL() time.Duration { return defaultTTL }

func CaptchaTTL() time.Duration { return captchaTTL }

func Key(parts ...string) string {
	return keyPrefix + strings.Join(parts, ":")
}

func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func GetString(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
}

func GetDelString(ctx context.Context, key string) (string, error) {
	return client.GetDel(ctx, key).Result()
}

func Del(ctx context.Context, keys ...string) error {
	return client.Del(ctx, keys...).Err()
}

func Exists(ctx context.Context, key string) (bool, error) {
	n, err := client.Exists(ctx, key).Result()
	return n > 0, err
}

func SetNX(ctx context.Context, key string, value any, ttl time.Duration) (bool, error) {
	return client.SetNX(ctx, key, value, ttl).Result()
}

func RotateValue(ctx context.Context, key, oldValue, newValue string, ttl time.Duration) (bool, error) {
	result, err := rotateValueScript.Run(
		ctx,
		client,
		[]string{key},
		oldValue,
		newValue,
		ttl.Milliseconds(),
	).Int()
	return result == 1, err
}

func SetJSON(ctx context.Context, key string, value any, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Set(ctx, key, data, ttl)
}

func GetJSON(ctx context.Context, key string, value any) error {
	data, err := client.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func DeleteByPattern(ctx context.Context, pattern string) error {
	var cursor uint64
	for {
		keys, next, err := client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := client.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = next
		if cursor == 0 {
			return nil
		}
	}
}

func IsNotFound(err error) bool {
	return err == redis.Nil
}

// RefreshTokenKey 刷新token的Redis key
func RefreshTokenKey(userID uint, sessionID string) string {
	return Key("auth", "refresh", fmt.Sprint(userID), sessionID)
}

// BlacklistKey 已注销token的黑名单key
func BlacklistKey(tokenID string) string {
	return Key("auth", "blacklist", tokenID)
}

// CaptchaKey 验证码key
func CaptchaKey(captchaID string) string {
	return Key("captcha", captchaID)
}

func UserMenusKey(userID uint) string {
	return Key("user", "menus", fmt.Sprint(userID))
}

func UserPermissionsKey(userID uint) string {
	return Key("user", "permissions", fmt.Sprint(userID))
}

func UserAuthorizationPattern() string {
	return Key("user", "*")
}

func RateLimitKey(scope, identifier string) string {
	return Key("rate_limit", scope, identifier)
}

func durationOrDefault(seconds int64, fallback time.Duration) time.Duration {
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
