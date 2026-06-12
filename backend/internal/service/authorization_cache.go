package service

import (
	"context"

	"go.uber.org/zap"
	"zhanxu-admin/backend/pkg/cache"
	"zhanxu-admin/backend/pkg/logger"
)

func getAuthorizationCache(key string, value any) bool {
	if err := cache.GetJSON(context.Background(), key, value); err != nil {
		if !cache.IsNotFound(err) {
			logger.Warn("读取权限缓存失败，回源数据库", zap.String("key", key), zap.Error(err))
		}
		return false
	}
	return true
}

func setAuthorizationCache(key string, value any) {
	if err := cache.SetJSON(context.Background(), key, value, cache.DefaultTTL()); err != nil {
		logger.Warn("写入权限缓存失败", zap.String("key", key), zap.Error(err))
	}
}

func invalidateUserAuthorizationCache(userID uint) {
	if err := cache.Del(
		context.Background(),
		cache.UserMenusKey(userID),
		cache.UserPermissionsKey(userID),
	); err != nil {
		logger.Warn("清除用户权限缓存失败", zap.Uint("user_id", userID), zap.Error(err))
	}
}

func invalidateAllAuthorizationCache() {
	if err := cache.DeleteByPattern(context.Background(), cache.UserAuthorizationPattern()); err != nil {
		logger.Warn("清除全部权限缓存失败", zap.Error(err))
	}
}
