package service

import (
	"errors"
	"fmt"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"zhanxu-admin/backend/internal/dto"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/response"
)

type RoleService struct {
	roleRepo *repository.RoleRepo
	enforcer *casbin.Enforcer
}

func NewRoleService(roleRepo *repository.RoleRepo, enforcer *casbin.Enforcer) *RoleService {
	return &RoleService{roleRepo: roleRepo, enforcer: enforcer}
}

func (s *RoleService) List(req *dto.RoleListReq) ([]dto.RoleResp, int64, error) {
	where := map[string]any{}
	if req.Name != "" {
		where["name LIKE ?"] = "%" + req.Name + "%"
	}
	if req.Code != "" {
		where["code LIKE ?"] = "%" + req.Code + "%"
	}
	if req.Status != nil {
		where["status = ?"] = *req.Status
	}

	roles, total, err := s.roleRepo.List(req.Page, req.PageSize, where)
	if err != nil {
		return nil, 0, err
	}
	resp := make([]dto.RoleResp, len(roles))
	for i, r := range roles {
		resp[i] = toRoleResp(&r)
	}
	return resp, total, nil
}

func (s *RoleService) All() ([]dto.RoleResp, error) {
	roles, err := s.roleRepo.FindAll()
	if err != nil {
		return nil, err
	}
	resp := make([]dto.RoleResp, len(roles))
	for i, r := range roles {
		resp[i] = toRoleResp(&r)
	}
	return resp, nil
}

func (s *RoleService) Get(id uint) (*dto.RoleResp, error) {
	r, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &BizError{Code: response.CodeRoleNotFound}
		}
		return nil, err
	}
	resp := toRoleResp(r)
	return &resp, nil
}

func (s *RoleService) Create(req *dto.CreateRoleReq) error {
	exists, err := s.roleRepo.ExistsByCode(req.Code)
	if err != nil {
		return err
	}
	if exists {
		return &BizError{Code: response.CodeRoleCodeExists}
	}
	role := &model.SysRole{
		Name:   req.Name,
		Code:   req.Code,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}
	if role.Status == 0 {
		role.Status = 1
	}
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(id uint, req *dto.UpdateRoleReq) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeRoleNotFound}
		}
		return err
	}
	exists, err := s.roleRepo.ExistsByCode(req.Code, id)
	if err != nil {
		return err
	}
	if exists {
		return &BizError{Code: response.CodeRoleCodeExists}
	}
	role.Name = req.Name
	role.Code = req.Code
	role.Sort = req.Sort
	role.Status = req.Status
	role.Remark = req.Remark
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	if _, err := s.roleRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeRoleNotFound}
		}
		return err
	}
	hasUsers, err := s.roleRepo.HasUsers(id)
	if err != nil {
		return err
	}
	if hasUsers {
		return &BizError{Code: response.CodeRoleInUse}
	}
	return s.roleRepo.Delete(id)
}

func (s *RoleService) UpdateStatus(id uint, req *dto.UpdateStatusReq) error {
	role, err := s.roleRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeRoleNotFound}
		}
		return err
	}
	role.Status = req.Status
	return s.roleRepo.Update(role)
}

func (s *RoleService) AssignMenus(id uint, req *dto.AssignMenusReq) error {
	if _, err := s.roleRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeRoleNotFound}
		}
		return err
	}
	if err := s.roleRepo.AssignMenus(id, req.MenuIDs); err != nil {
		return err
	}
	// 同步更新 Casbin 权限策略（此处简化，实际按菜单permission字段同步接口权限）
	_ = s.enforcer.LoadPolicy()
	return nil
}

func (s *RoleService) GetMenuIDs(id uint) ([]uint, error) {
	return s.roleRepo.GetMenuIDs(id)
}

func toRoleResp(r *model.SysRole) dto.RoleResp {
	return dto.RoleResp{
		ID:        r.ID,
		Name:      r.Name,
		Code:      r.Code,
		Sort:      r.Sort,
		Status:    r.Status,
		Remark:    r.Remark,
		CreatedAt: r.CreatedAt,
	}
}

func roleSub(roleCode string) string {
	return fmt.Sprintf("role:%s", roleCode)
}
