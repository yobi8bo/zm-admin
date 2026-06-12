package repository

import (
	"gorm.io/gorm"
	"zhanxu-admin/backend/internal/model"
)

type RoleRepo struct{ db *gorm.DB }

func NewRoleRepo(db *gorm.DB) *RoleRepo { return &RoleRepo{db: db} }

func (r *RoleRepo) FindByID(id uint) (*model.SysRole, error) {
	var role model.SysRole
	err := r.db.First(&role, id).Error
	return &role, err
}

func (r *RoleRepo) FindAll() ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Where("status = 1").Order("sort ASC").Find(&roles).Error
	return roles, err
}

func (r *RoleRepo) List(page, pageSize int, where map[string]any) ([]model.SysRole, int64, error) {
	var roles []model.SysRole
	var total int64
	q := r.db.Model(&model.SysRole{})
	for k, v := range where {
		q = q.Where(k, v)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("sort ASC").Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepo) Create(role *model.SysRole) error {
	return r.db.Create(role).Error
}

func (r *RoleRepo) Update(role *model.SysRole) error {
	return r.db.Save(role).Error
}

func (r *RoleRepo) Delete(id uint) error {
	return r.db.Delete(&model.SysRole{}, id).Error
}

func (r *RoleRepo) ExistsByCode(code string, excludeID ...uint) (bool, error) {
	q := r.db.Model(&model.SysRole{}).Where("code = ?", code)
	if len(excludeID) > 0 && excludeID[0] > 0 {
		q = q.Where("id != ?", excludeID[0])
	}
	var count int64
	return count > 0, q.Count(&count).Error
}

func (r *RoleRepo) HasUsers(id uint) (bool, error) {
	var count int64
	err := r.db.Table("sys_user_role").Where("role_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *RoleRepo) AssignMenus(roleID uint, menuIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var role model.SysRole
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}
		var menus []model.SysMenu
		if err := tx.Find(&menus, menuIDs).Error; err != nil {
			return err
		}
		return tx.Model(&role).Association("Menus").Replace(menus)
	})
}

func (r *RoleRepo) GetMenuIDs(roleID uint) ([]uint, error) {
	var role model.SysRole
	if err := r.db.Preload("Menus").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, len(role.Menus))
	for i, m := range role.Menus {
		ids[i] = m.ID
	}
	return ids, nil
}

func (r *RoleRepo) FindByIDs(ids []uint) ([]model.SysRole, error) {
	var roles []model.SysRole
	err := r.db.Find(&roles, ids).Error
	return roles, err
}
