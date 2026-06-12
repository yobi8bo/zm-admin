package service

import (
	"errors"

	"gorm.io/gorm"
	"zhanxu-admin/backend/internal/dto"
	"zhanxu-admin/backend/internal/model"
	"zhanxu-admin/backend/internal/repository"
	"zhanxu-admin/backend/pkg/response"
)

type MenuService struct {
	menuRepo *repository.MenuRepo
}

func NewMenuService(menuRepo *repository.MenuRepo) *MenuService {
	return &MenuService{menuRepo: menuRepo}
}

func (s *MenuService) List() ([]dto.MenuResp, error) {
	menus, err := s.menuRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return buildMenuTree(menus, 0), nil
}

func (s *MenuService) Get(id uint) (*dto.MenuResp, error) {
	m, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &BizError{Code: response.CodeMenuNotFound}
		}
		return nil, err
	}
	r := toMenuResp(m)
	return &r, nil
}

func (s *MenuService) Create(req *dto.CreateMenuReq) error {
	m := &model.SysMenu{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Type:       req.Type,
		Path:       req.Path,
		Component:  req.Component,
		Permission: req.Permission,
		Icon:       req.Icon,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
	}
	if m.Status == 0 {
		m.Status = 1
	}
	if m.Visible == 0 {
		m.Visible = 1
	}
	if err := s.menuRepo.Create(m); err != nil {
		return err
	}
	invalidateAllAuthorizationCache()
	return nil
}

func (s *MenuService) Update(id uint, req *dto.UpdateMenuReq) error {
	m, err := s.menuRepo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeMenuNotFound}
		}
		return err
	}
	m.ParentID = req.ParentID
	m.Name = req.Name
	m.Type = req.Type
	m.Path = req.Path
	m.Component = req.Component
	m.Permission = req.Permission
	m.Icon = req.Icon
	m.Sort = req.Sort
	m.Visible = req.Visible
	m.Status = req.Status
	if err := s.menuRepo.Update(m); err != nil {
		return err
	}
	invalidateAllAuthorizationCache()
	return nil
}

func (s *MenuService) Delete(id uint) error {
	if _, err := s.menuRepo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &BizError{Code: response.CodeMenuNotFound}
		}
		return err
	}
	hasChildren, err := s.menuRepo.HasChildren(id)
	if err != nil {
		return err
	}
	if hasChildren {
		return &BizError{Code: response.CodeMenuHasChildren}
	}
	if err := s.menuRepo.Delete(id); err != nil {
		return err
	}
	invalidateAllAuthorizationCache()
	return nil
}

func toMenuResp(m *model.SysMenu) dto.MenuResp {
	return dto.MenuResp{
		ID:         m.ID,
		ParentID:   m.ParentID,
		Name:       m.Name,
		Type:       m.Type,
		Path:       m.Path,
		Component:  m.Component,
		Permission: m.Permission,
		Icon:       m.Icon,
		Sort:       m.Sort,
		Visible:    m.Visible,
		Status:     m.Status,
	}
}

func buildMenuTree(menus []model.SysMenu, parentID uint) []dto.MenuResp {
	var nodes []dto.MenuResp
	for _, m := range menus {
		if m.ParentID == parentID {
			node := toMenuResp(&m)
			node.Children = buildMenuTree(menus, m.ID)
			nodes = append(nodes, node)
		}
	}
	return nodes
}
