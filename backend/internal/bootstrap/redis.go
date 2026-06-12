package bootstrap

import (
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/pkg/cache"
	"zhanxu-admin/backend/pkg/logger"
)

func InitRedis(cfg config.RedisConfig) error {
	if err := cache.Init(cfg); err != nil {
		return err
	}
	logger.Info("Redis连接成功")
	return nil
}
