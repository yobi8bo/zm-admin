package repository

import (
	"gorm.io/gorm"
	"zm-project/backend/internal/model"
)

type DeptRepo struct{ db *gorm.DB }

func NewDeptRepo(db *gorm.DB) *DeptRepo { return &DeptRepo{db: db} }

func (r *DeptRepo) FindByID(id uint) (*model.SysDept, error) {
	var d model.SysDept
	err := r.db.First(&d, id).Error
	return &d, err
}

func (r *DeptRepo) FindAll() ([]model.SysDept, error) {
	var depts []model.SysDept
	err := r.db.Order("sort ASC").Find(&depts).Error
	return depts, err
}

func (r *DeptRepo) Create(d *model.SysDept) error {
	return r.db.Create(d).Error
}

func (r *DeptRepo) Update(d *model.SysDept) error {
	return r.db.Save(d).Error
}

func (r *DeptRepo) Delete(id uint) error {
	return r.db.Delete(&model.SysDept{}, id).Error
}

func (r *DeptRepo) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SysDept{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *DeptRepo) HasUsers(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.SysUser{}).Where("dept_id = ?", id).Count(&count).Error
	return count > 0, err
}
