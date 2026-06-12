package service

import (
	"reflect"
	"testing"

	"zhanxu-admin/backend/internal/model"
)

func TestToUserRespIncludesRoleIDs(t *testing.T) {
	user := &model.SysUser{
		Roles: []model.SysRole{
			{Base: model.Base{ID: 2}},
			{Base: model.Base{ID: 5}},
		},
	}

	resp := toUserResp(user)

	if want := []uint{2, 5}; !reflect.DeepEqual(resp.RoleIDs, want) {
		t.Fatalf("RoleIDs = %v, want %v", resp.RoleIDs, want)
	}
}
