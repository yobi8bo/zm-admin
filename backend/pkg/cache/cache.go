package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"zhanxu-admin/backend/config"
)

var client *redis.Client

func Init(cfg config.RedisConfig) error {
	client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return client.Ping(ctx).Err()
}

func Get() *redis.Client { return client }

func Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	return client.Set(ctx, key, value, ttl).Err()
}

func GetString(ctx context.Context, key string) (string, error) {
	return client.Get(ctx, key).Result()
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

// RefreshTokenKey 刷新token的Redis key
func RefreshTokenKey(userID uint) string {
	return fmt.Sprintf("refresh_token:%d", userID)
}

// BlacklistKey 已注销token的黑名单key
func BlacklistKey(token string) string {
	return fmt.Sprintf("token_blacklist:%s", token)
}

// CaptchaKey 验证码key
func CaptchaKey(captchaID string) string {
	return fmt.Sprintf("captcha:%s", captchaID)
}
