package service

import (
	"errors"
	"time"

	"github.com/casbin/casbin/v2"
	"gorm.io/gorm"
	"zhanxu-admin/backend/internal/dto"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/crypto"
	"zhanxu-admin/backend/pkg/response"
)

type UserService struct {
	userRepo *repository.UserRepo
	roleRepo *repository.RoleRepo
	menuRepo *repository.MenuRepo
	enforcer *casbin.Enforcer
}

func NewUserService(userRepo *repository.UserRepo, roleRepo *repository.RoleRepo,
	menuRepo *repository.MenuRepo, enforcer *casbin.Enforcer) *UserService {
	return &UserService{userRepo: userRepo, roleRepo: roleRepo, menuRepo: menuRepo, enforcer: enforcer}
}

func (s *UserService) List(req *dto.UserListReq) ([]dto.UserResp, int64, error) {
	where := map[string]any{}
	if req.Username != "" {
		where["username LIKE ?"] = "%" + req.Username + "%"
	}
	if req.Phone != "" {
		where["phone = ?"] = req.Phone
	}
	if req.Status != nil {
		where["status = ?"] = *req.Status
	}
	if req.DeptID > 0 {
		where["dept_id = ?"] = req.DeptID
	}

	users, total, err := s.userRepo.List(req.Page, req.PageSize, where)
	if err != nil {
		return nil, 0, err
	}

	resp := make([]dto.UserResp, len(users))
	for i, u := range users {
		resp[i] = toUserResp(&u)
	}
	return resp, total, nil
}

func (s *UserService) Get(id uint) (*dto.UserResp, error) {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &BizError{Code: response.CodeUserNotFound}
		}
		return nil, err
	}
	r := toUserResp(u)
	return &r, nil
}

func (s *UserService) Create(req *dto.CreateUserReq) error {
	exists, err := s.userRepo.ExistsByUsername(req.Username)
	if err != nil {
		return err
	}
	if exists {
		return &BizError{Code: response.CodeUsernameExists}
	}

	hashed, err := crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}

	u := &model.SysUser{
		DeptID:   req.DeptID,
		Username: req.Username,
		Nickname: req.Nickname,
		Password: hashed,
		Email:    req.Email,
		Phone:    req.Phone,
		Gender:   req.Gender,
		Status:   req.Status,
		Remark:   req.Remark,
	}
	if u.Status == 0 {
		u.Status = 1
	}
	return s.userRepo.Create(u, req.RoleIDs)
}

func (s *UserService) Update(id uint, req *dto.UpdateUserReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeUserNotFound}
		}
		return err
	}
	u.DeptID = req.DeptID
	u.Nickname = req.Nickname
	u.Email = req.Email
	u.Phone = req.Phone
	u.Gender = req.Gender
	u.Remark = req.Remark
	return s.userRepo.Update(u)
}

func (s *UserService) Delete(id uint) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeUserNotFound}
		}
		return err
	}
	if isAdminUser(u) {
		return &BizError{Code: response.CodeAdminUserProtected}
	}
	return s.userRepo.Delete(id)
}

func isAdminUser(u *model.SysUser) bool {
	if u.Username == "admin" {
		return true
	}
	for _, role := range u.Roles {
		if role.Code == "admin" {
			return true
		}
	}
	return false
}

func (s *UserService) UpdateStatus(id uint, req *dto.UpdateStatusReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeUserNotFound}
		}
		return err
	}
	u.Status = req.Status
	return s.userRepo.Update(u)
}

func (s *UserService) ResetPassword(id uint, req *dto.ResetPasswordReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeUserNotFound}
		}
		return err
	}
	hashed, err := crypto.HashPassword(req.Password)
	if err != nil {
		return err
	}
	u.Password = hashed
	return s.userRepo.Update(u)
}

func (s *UserService) AssignRoles(id uint, req *dto.AssignRolesReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeUserNotFound}
		}
		return err
	}
	if isAdminUser(u) {
		return &BizError{Code: response.CodeAdminRoleProtected}
	}
	return s.userRepo.AssignRoles(id, req.RoleIDs)
}

func (s *UserService) GetMe(id uint) (*dto.UserResp, error) {
	return s.Get(id)
}

func (s *UserService) UpdateMe(id uint, req *dto.UpdateMeReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if req.Nickname != "" {
		u.Nickname = req.Nickname
	}
	u.Email = req.Email
	u.Phone = req.Phone
	u.Gender = req.Gender
	if req.Avatar != "" {
		u.Avatar = req.Avatar
	}
	return s.userRepo.Update(u)
}

func (s *UserService) UpdateMyPassword(id uint, req *dto.UpdateMyPasswordReq) error {
	u, err := s.userRepo.FindByID(id)
	if err != nil {
		return err
	}
	if !crypto.CheckPassword(req.OldPassword, u.Password) {
		return &BizError{Code: response.CodePasswordError}
	}
	hashed, err := crypto.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	u.Password = hashed
	return s.userRepo.Update(u)
}

func (s *UserService) GetMyMenus(userID uint) ([]dto.RouteResp, error) {
	roleIDs, err := s.userRepo.GetRoleIDs(userID)
	if err != nil {
		return nil, err
	}
	menus, err := s.menuRepo.FindByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	// 过滤掉按钮类型，只返回目录和菜单
	var routes []model.SysMenu
	for _, m := range menus {
		if m.Type != 3 {
			routes = append(routes, m)
		}
	}
	return buildRouteTree(routes, 0), nil
}

func (s *UserService) GetMyPermissions(userID uint) ([]string, error) {
	roleIDs, err := s.userRepo.GetRoleIDs(userID)
	if err != nil {
		return nil, err
	}
	menus, err := s.menuRepo.FindByRoleIDs(roleIDs)
	if err != nil {
		return nil, err
	}
	var perms []string
	for _, m := range menus {
		if m.Type == 3 && m.Permission != "" {
			perms = append(perms, m.Permission)
		}
	}
	return perms, nil
}

func toUserResp(u *model.SysUser) dto.UserResp {
	r := dto.UserResp{
		ID:        u.ID,
		DeptID:    u.DeptID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Email:     u.Email,
		Phone:     u.Phone,
		Gender:    u.Gender,
		Status:    u.Status,
		IsAdmin:   isAdminUser(u),
		RoleIDs:   make([]uint, len(u.Roles)),
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
	}
	for i, role := range u.Roles {
		r.RoleIDs[i] = role.ID
	}
	if u.Dept != nil {
		r.DeptName = u.Dept.Name
	}
	return r
}

func buildRouteTree(menus []model.SysMenu, parentID uint) []dto.RouteResp {
	var nodes []dto.RouteResp
	for _, m := range menus {
		if m.ParentID == parentID {
			node := dto.RouteResp{
				ID:        m.ID,
				Name:      m.Name,
				Path:      m.Path,
				Component: m.Component,
				Icon:      m.Icon,
				Sort:      m.Sort,
				Children:  buildRouteTree(menus, m.ID),
			}
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func toTime(t time.Time) time.Time { return t }
