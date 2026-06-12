package repository

import (
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"zhanxu-admin/backend/internal/model"
)

type UserRepo struct{ db *gorm.DB }

var ErrUsernameExists = errors.New("username already exists")

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) FindByID(id uint) (*model.SysUser, error) {
	var u model.SysUser
	err := r.db.Preload("Dept").Preload("Roles").First(&u, id).Error
	return &u, err
}

func (r *UserRepo) FindByUsername(username string) (*model.SysUser, error) {
	var u model.SysUser
	err := r.db.Where("username = ?", username).First(&u).Error
	return &u, err
}

func (r *UserRepo) List(page, pageSize int, where map[string]any) ([]model.SysUser, int64, error) {
	var users []model.SysUser
	var total int64
	q := r.db.Model(&model.SysUser{}).Preload("Dept").Preload("Roles")
	for k, v := range where {
		q = q.Where(k, v)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := q.Offset((page - 1) * pageSize).Limit(pageSize).
		Order("created_at DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepo) Create(u *model.SysUser, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var deleted model.SysUser
		err := tx.Unscoped().
			Clauses(clause.Locking{Strength: "UPDATE"}).
			Where("username = ? AND deleted_at IS NOT NULL", u.Username).
			First(&deleted).Error
		switch {
		case err == nil:
			if err := restoreDeletedUser(tx, &deleted, u); err != nil {
				return err
			}
			u.ID = deleted.ID
		case errors.Is(err, gorm.ErrRecordNotFound):
			if err := tx.Create(u).Error; err != nil {
				if isDuplicateEntry(err) {
					return ErrUsernameExists
				}
				return err
			}
		default:
			return err
		}

		var roles []model.SysRole
		if len(roleIDs) > 0 {
			if err := tx.Find(&roles, roleIDs).Error; err != nil {
				return err
			}
		}
		return tx.Model(u).Association("Roles").Replace(roles)
	})
}

func restoreDeletedUser(tx *gorm.DB, deleted, replacement *model.SysUser) error {
	return tx.Unscoped().Model(deleted).Updates(map[string]any{
		"deleted_at": nil,
		"dept_id":    replacement.DeptID,
		"nickname":   replacement.Nickname,
		"password":   replacement.Password,
		"avatar":     replacement.Avatar,
		"email":      replacement.Email,
		"phone":      replacement.Phone,
		"gender":     replacement.Gender,
		"status":     replacement.Status,
		"last_login": nil,
		"remark":     replacement.Remark,
	}).Error
}

func isDuplicateEntry(err error) bool {
	var mysqlErr *mysql.MySQLError
	return errors.As(err, &mysqlErr) && mysqlErr.Number == 1062
}

func (r *UserRepo) Update(u *model.SysUser) error {
	return r.db.Save(u).Error
}

func (r *UserRepo) Delete(id uint) error {
	return r.db.Delete(&model.SysUser{}, id).Error
}

func (r *UserRepo) AssignRoles(userID uint, roleIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var user model.SysUser
		if err := tx.First(&user, userID).Error; err != nil {
			return err
		}
		var roles []model.SysRole
		if err := tx.Find(&roles, roleIDs).Error; err != nil {
			return err
		}
		return tx.Model(&user).Association("Roles").Replace(roles)
	})
}

func (r *UserRepo) ExistsByUsername(username string, excludeID ...uint) (bool, error) {
	q := r.db.Model(&model.SysUser{}).Where("username = ?", username)
	if len(excludeID) > 0 && excludeID[0] > 0 {
		q = q.Where("id != ?", excludeID[0])
	}
	var count int64
	return count > 0, q.Count(&count).Error
}

func (r *UserRepo) UpdateLastLogin(id uint) error {
	return r.db.Model(&model.SysUser{}).Where("id = ?", id).
		Update("last_login", gorm.Expr("NOW()")).Error
}

func (r *UserRepo) GetRoleIDs(userID uint) ([]uint, error) {
	var user model.SysUser
	if err := r.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	ids := make([]uint, len(user.Roles))
	ids = ids[:0]
	for _, role := range user.Roles {
		if role.Status == 1 {
			ids = append(ids, role.ID)
		}
	}
	return ids, nil
}
