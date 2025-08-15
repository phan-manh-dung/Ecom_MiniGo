package repository

// repository là nơi lưu trữ dữ liệu vào database

import (
	"gin/shared/generic"
	"gin/user_service/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface định nghĩa interface cho UserRepository
type UserRepositoryInterface interface {
	GetByID(id uint) (*model.User, error)
	Create(user *model.User) error
	CreateAccount(account *model.Account) error
	Update(user *model.User) error
	Delete(id uint) error
	GetAll(page, limit int) ([]model.User, int64, error)
	GetBySDT(sdt string) (*model.User, error)
	GetRoleByID(id uint) (*model.Role, error)
	ListRoles(page, limit int) ([]model.Role, int64, error)
}

type UserRepository struct {
	*generic.BaseRepository
}

// Đảm bảo UserRepository implement UserRepositoryInterface
var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{
		BaseRepository: generic.NewBaseRepository(db),
	}
}

// Lấy user theo ID với preload Account.Role
func (r *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	result := r.GetDB().Preload("Account.Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Tạo user
func (r *UserRepository) Create(user *model.User) error {
	return generic.GenericCreate(r.GetDB(), user)
}

// Tạo account
func (r *UserRepository) CreateAccount(account *model.Account) error {
	return generic.GenericCreate(r.GetDB(), account)
}

// Cập nhật user
func (r *UserRepository) Update(user *model.User) error {
	return generic.GenericUpdate(r.GetDB(), user)
}

// Xóa user
func (r *UserRepository) Delete(id uint) error {
	return generic.GenericDelete[model.User, uint](r.GetDB(), id)
}

// Phân trang với offset/limit
func (r *UserRepository) GetAll(page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Count total
	if err := r.GetDB().Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	result := r.GetDB().Preload("Account.Role").Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// Tìm user theo số điện thoại với preload Account.Role
func (r *UserRepository) GetBySDT(sdt string) (*model.User, error) {
	var user model.User
	result := r.GetDB().Preload("Account.Role").Where("sdt = ?", sdt).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetRoleByID(id uint) (*model.Role, error) {
	return generic.GenericGetByID[model.Role, uint](r.GetDB(), id)
}

func (r *UserRepository) ListRoles(page, limit int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64
	if err := r.GetDB().Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	result := r.GetDB().Offset(offset).Limit(limit).Find(&roles)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return roles, total, nil
}
