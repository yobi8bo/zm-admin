package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/config"
	"zhanxu-admin/backend/pkg/cache"
	"zhanxu-admin/backend/pkg/response"
)

func RateLimit(cfg config.RateLimitConfig) gin.HandlerFunc {
	return RateLimitByIP("global", cfg.Rate, cfg.Burst)
}

func RateLimitByIP(scope string, rate float64, burst int) gin.HandlerFunc {
	return func(c *gin.Context) {
		allowed, err := cache.Allow(context.Background(), cache.RateLimitKey(scope, c.ClientIP()), rate, burst)
		if err != nil {
			response.ServerError(c)
			c.Abort()
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusOK, gin.H{
				"code": response.CodeTooManyReqs,
				"msg":  response.GetMsg(response.CodeTooManyReqs),
				"data": nil,
			})
			return
		}
		c.Next()
	}
}
