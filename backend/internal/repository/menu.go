package repository

import (
	"gorm.io/gorm"
	"zm-project/backend/internal/model"
)

type MenuRepo struct{ db *gorm.DB }

func NewMenuRepo(db *gorm.DB) *MenuRepo { return &MenuRepo{db: db} }

func (r *MenuRepo) FindByID(id uint) (*model.SysMenu, error) {
	var m model.SysMenu
	err := r.db.First(&m, id).Error
	return &m, err
}

func (r *MenuRepo) FindAll() ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (r *MenuRepo) FindByStatus(status int8) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.Where("status = ?", status).Order("sort ASC").Find(&menus).Error
	return menus, err
}

func (r *MenuRepo) Create(m *model.SysMenu) error {
	return r.db.Create(m).Error
}

func (r *MenuRepo) Update(m *model.SysMenu) error {
	return r.db.Save(m).Error
}

func (r *MenuRepo) Delete(id uint) error {
	return r.db.Delete(&model.SysMenu{}, id).Error
}

func (r *MenuRepo) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SysMenu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

// FindByRoleIDs 根据角色ID列表查询有权限的菜单
func (r *MenuRepo) FindByRoleIDs(roleIDs []uint) ([]model.SysMenu, error) {
	var menus []model.SysMenu
	err := r.db.
		Joins("JOIN sys_role_menu ON sys_role_menu.sys_menu_id = sys_menu.id").
		Where("sys_role_menu.sys_role_id IN ? AND sys_menu.status = 1", roleIDs).
		Order("sys_menu.sort ASC").
		Distinct().
		Find(&menus).Error
	return menus, err
}
