package repository

import (
	"gorm.io/gorm"
	"zm-project/backend/internal/model"
)

type LogRepo struct{ db *gorm.DB }

func NewLogRepo(db *gorm.DB) *LogRepo { return &LogRepo{db: db} }

func (r *LogRepo) CreateOperationLog(log *model.SysOperationLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepo) ListOperationLog(page, pageSize int, where map[string]any) ([]model.SysOperationLog, int64, error) {
	var logs []model.SysOperationLog
	var total int64
	q := r.db.Model(&model.SysOperationLog{})
	for k, v := range where {
		q = q.Where(k, v)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepo) ClearOperationLog() error {
	return r.db.Where("1 = 1").Delete(&model.SysOperationLog{}).Error
}

func (r *LogRepo) CreateLoginLog(log *model.SysLoginLog) error {
	return r.db.Create(log).Error
}

func (r *LogRepo) ListLoginLog(page, pageSize int, where map[string]any) ([]model.SysLoginLog, int64, error) {
	var logs []model.SysLoginLog
	var total int64
	q := r.db.Model(&model.SysLoginLog{})
	for k, v := range where {
		q = q.Where(k, v)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepo) ClearLoginLog() error {
	return r.db.Where("1 = 1").Delete(&model.SysLoginLog{}).Error
}
