package repository

import (
	"gin/user_service/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Account.Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) GetAll(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Count total
	if err := r.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	result := r.db.Preload("Account.Role").Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

func (r *UserRepository) GetBySDT(sdt string) (*model.User, error) {
	var user model.User
	result := r.db.Preload("Account.Role").Where("sdt = ?", sdt).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetRoleByID(id uint) (*model.Role, error) {
	var role model.Role
	result := r.db.First(&role, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &role, nil
}

func (r *UserRepository) ListRoles(page, limit int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64
	if err := r.db.Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	result := r.db.Offset(offset).Limit(limit).Find(&roles)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return roles, total, nil
}
