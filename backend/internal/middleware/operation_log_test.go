package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"zhanxu-admin/backend/pkg/response"
)

func TestShouldLogOperation(t *testing.T) {
	tests := []struct {
		method string
		path   string
		want   bool
	}{
		{http.MethodGet, "/api/v1/users", false},
		{http.MethodPost, "/api/v1/users", true},
		{http.MethodPut, "/api/v1/users/2", true},
		{http.MethodDelete, "/api/v1/users/2", true},
		{http.MethodDelete, "/api/v1/logs/operation", false},
	}

	for _, tt := range tests {
		if got := shouldLogOperation(tt.method, tt.path); got != tt.want {
			t.Fatalf("shouldLogOperation(%q, %q) = %v, want %v", tt.method, tt.path, got, tt.want)
		}
	}
}

func TestRedactSensitiveFields(t *testing.T) {
	value := map[string]any{
		"username": "tester",
		"password": "secret",
		"nested": map[string]any{
			"old_password": "old-secret",
			"nickname":     "测试用户",
		},
	}

	redactSensitiveFields(value)

	if value["password"] != "***" {
		t.Fatalf("password was not redacted")
	}
	nested := value["nested"].(map[string]any)
	if nested["old_password"] != "***" {
		t.Fatalf("nested password was not redacted")
	}
	if nested["nickname"] != "测试用户" {
		t.Fatalf("non-sensitive field was changed")
	}
}

func TestOperationDescription(t *testing.T) {
	tests := []struct {
		route      string
		method     string
		wantModule string
		wantAction string
	}{
		{"/api/v1/users", http.MethodPost, "用户管理", "新增"},
		{"/api/v1/users/:id/roles", http.MethodPut, "用户管理", "分配角色"},
		{"/api/v1/roles/:id/menus", http.MethodPut, "角色管理", "分配菜单"},
		{"/api/v1/depts/:id", http.MethodDelete, "部门管理", "删除"},
	}

	for _, tt := range tests {
		module, action := operationDescription(tt.route, tt.method)
		if module != tt.wantModule || action != tt.wantAction {
			t.Fatalf("operationDescription(%q, %q) = (%q, %q), want (%q, %q)",
				tt.route, tt.method, module, action, tt.wantModule, tt.wantAction)
		}
	}
}

func TestOperationResult(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name       string
		body       string
		contextErr error
		httpStatus int
		wantStatus int
		wantError  string
	}{
		{
			name:       "business success",
			body:       `{"code":200,"msg":"success","data":null}`,
			httpStatus: http.StatusOK,
			wantStatus: 200,
		},
		{
			name:       "business failure",
			body:       `{"code":10004,"msg":"用户名已存在","data":null}`,
			httpStatus: http.StatusOK,
			wantStatus: 10004,
			wantError:  "用户名已存在",
		},
		{
			name:       "context error",
			contextErr: errors.New("unexpected failure"),
			httpStatus: http.StatusOK,
			wantStatus: 500,
			wantError:  "Error #01: unexpected failure\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			recorder.WriteHeader(tt.httpStatus)
			c, _ := gin.CreateTestContext(recorder)
			if tt.contextErr != nil {
				_ = c.Error(tt.contextErr)
			}

			status, errMsg := operationResult([]byte(tt.body), c)
			if status != tt.wantStatus || errMsg != tt.wantError {
				t.Fatalf("operationResult() = (%d, %q), want (%d, %q)",
					status, errMsg, tt.wantStatus, tt.wantError)
			}
		})
	}
}

func TestOperationLogStatusFallsBackToFailureWhenErrorExists(t *testing.T) {
	status := normalizeOperationStatus(response.CodeSuccess, "operation failed")
	if status != response.CodeServerError {
		t.Fatalf("status = %d, want %d", status, response.CodeServerError)
	}
}
