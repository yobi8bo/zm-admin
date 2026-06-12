package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"zm-project/backend/pkg/cache"
	"zm-project/backend/pkg/jwtutil"
	"zm-project/backend/pkg/response"
)

const (
	CtxUserID   = "userID"
	CtxUsername = "username"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := extractToken(c)
		if token == "" {
			response.Fail(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		// 检查黑名单
		blacklisted, _ := cache.Exists(context.Background(), cache.BlacklistKey(token))
		if blacklisted {
			response.Fail(c, response.CodeTokenInvalid)
			c.Abort()
			return
		}

		claims, err := jwtutil.ParseToken(token, jwtSecret)
		if err != nil {
			if err == jwtutil.ErrTokenExpired {
				response.Fail(c, response.CodeTokenExpired)
			} else {
				response.Fail(c, response.CodeTokenInvalid)
			}
			c.Abort()
			return
		}

		c.Set(CtxUserID, claims.UserID)
		c.Set(CtxUsername, claims.Username)
		c.Next()
	}
}

func extractToken(c *gin.Context) string {
	bearer := c.GetHeader("Authorization")
	if strings.HasPrefix(bearer, "Bearer ") {
		return strings.TrimPrefix(bearer, "Bearer ")
	}
	return ""
}

func GetUserID(c *gin.Context) uint {
	id, _ := c.Get(CtxUserID)
	if uid, ok := id.(uint); ok {
		return uid
	}
	return 0
}

func GetUsername(c *gin.Context) string {
	name, _ := c.Get(CtxUsername)
	if username, ok := name.(string); ok {
		return username
	}
	return ""
}
