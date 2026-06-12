package service

import (
	"testing"

	"zhanxu-admin/backend/internal/model"
)

func TestIsAdminRole(t *testing.T) {
	tests := []struct {
		code string
		want bool
	}{
		{code: "admin", want: true},
		{code: "user", want: false},
		{code: "administrator", want: false},
	}

	for _, tt := range tests {
		role := &model.SysRole{Code: tt.code}
		if got := isAdminRole(role); got != tt.want {
			t.Fatalf("isAdminRole(%q) = %v, want %v", tt.code, got, tt.want)
		}
	}
}
