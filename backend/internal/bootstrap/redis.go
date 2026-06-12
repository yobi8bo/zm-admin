package bootstrap

import (
	"zm-project/backend/config"
	"zm-project/backend/pkg/cache"
	"zm-project/backend/pkg/logger"
)

func InitRedis(cfg config.RedisConfig) error {
	if err := cache.Init(cfg); err != nil {
		return err
	}
	logger.Info("Redis连接成功")
	return nil
}
