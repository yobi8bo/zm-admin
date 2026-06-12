package middleware

import (
	"net/http"
	"testing"
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
