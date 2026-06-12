package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/logger"
	"zhanxu-admin/backend/pkg/response"
)

const maxOperationLogBodySize = 64 * 1024

type operationResponseWriter struct {
	gin.ResponseWriter
	body bytes.Buffer
}

func (w *operationResponseWriter) Write(data []byte) (int, error) {
	if w.body.Len() < maxOperationLogBodySize {
		remaining := maxOperationLogBodySize - w.body.Len()
		w.body.Write(data[:min(len(data), remaining)])
	}
	return w.ResponseWriter.Write(data)
}

func OperationLog(logRepo *repository.LogRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !shouldLogOperation(c.Request.Method, c.Request.URL.Path) {
			c.Next()
			return
		}

		body := readAndRestoreBody(c)
		writer := &operationResponseWriter{ResponseWriter: c.Writer}
		c.Writer = writer
		start := time.Now()

		c.Next()

		module, action := operationDescription(c.FullPath(), c.Request.Method)
		logEntry := &model.SysOperationLog{
			UserID:    GetUserID(c),
			Username:  GetUsername(c),
			Module:    module,
			Action:    action,
			Method:    c.Request.Method,
			Path:      c.Request.URL.Path,
			Query:     c.Request.URL.RawQuery,
			Body:      body,
			IP:        c.ClientIP(),
			UserAgent: c.Request.UserAgent(),
			Status:    c.Writer.Status(),
			Latency:   time.Since(start).Milliseconds(),
			Error:     operationError(writer.body.Bytes(), c),
		}
		if err := logRepo.CreateOperationLog(logEntry); err != nil {
			logger.Error("写入操作日志失败", zap.Error(err))
		}
	}
}

func shouldLogOperation(method, path string) bool {
	if path == "/api/v1/logs/operation" {
		return false
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	default:
		return false
	}
}

func readAndRestoreBody(c *gin.Context) string {
	if c.Request.Body == nil || !strings.Contains(c.GetHeader("Content-Type"), "application/json") {
		return ""
	}
	data, err := io.ReadAll(io.LimitReader(c.Request.Body, maxOperationLogBodySize))
	if err != nil {
		return ""
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(data))

	var value any
	if json.Unmarshal(data, &value) != nil {
		return string(data)
	}
	redactSensitiveFields(value)
	sanitized, err := json.Marshal(value)
	if err != nil {
		return ""
	}
	return string(sanitized)
}

func redactSensitiveFields(value any) {
	switch typed := value.(type) {
	case map[string]any:
		for key, child := range typed {
			switch strings.ToLower(key) {
			case "password", "old_password", "new_password", "refresh_token", "access_token":
				typed[key] = "***"
			default:
				redactSensitiveFields(child)
			}
		}
	case []any:
		for _, child := range typed {
			redactSensitiveFields(child)
		}
	}
}

func operationError(data []byte, c *gin.Context) string {
	if len(c.Errors) > 0 {
		return c.Errors.String()
	}
	var resp response.Response
	if json.Unmarshal(data, &resp) == nil && resp.Code != response.CodeSuccess {
		return resp.Msg
	}
	return ""
}

func operationDescription(route, method string) (string, string) {
	module := "系统"
	switch {
	case strings.HasPrefix(route, "/api/v1/users"), strings.HasPrefix(route, "/api/v1/user/"):
		module = "用户管理"
	case strings.HasPrefix(route, "/api/v1/roles"):
		module = "角色管理"
	case strings.HasPrefix(route, "/api/v1/menus"):
		module = "菜单管理"
	case strings.HasPrefix(route, "/api/v1/depts"):
		module = "部门管理"
	case strings.HasPrefix(route, "/api/v1/logs"):
		module = "日志管理"
	case strings.HasPrefix(route, "/api/v1/files"):
		module = "文件管理"
	case strings.HasPrefix(route, "/api/v1/auth"):
		module = "认证管理"
	}

	action := map[string]string{
		http.MethodPost:   "新增",
		http.MethodPut:    "修改",
		http.MethodPatch:  "修改",
		http.MethodDelete: "删除",
	}[method]
	switch {
	case strings.HasSuffix(route, "/status"):
		action = "修改状态"
	case strings.HasSuffix(route, "/password"):
		action = "重置密码"
	case strings.HasSuffix(route, "/roles"):
		action = "分配角色"
	case strings.HasSuffix(route, "/menus"):
		action = "分配菜单"
	case route == "/api/v1/auth/logout":
		action = "退出登录"
	case route == "/api/v1/files/upload":
		action = "上传文件"
	}
	return module, action
}
