package repository

// repository là nơi lưu trữ dữ liệu vào database

import (
	"context"
	"gin/shared/generic"
	"gin/user_service/model"

	"gorm.io/gorm"
)

// UserRepositoryInterface định nghĩa interface cho UserRepository
type UserRepositoryInterface interface {
	GetByID(ctx context.Context, id uint) (*model.User, error)
	Create(ctx context.Context, user *model.User) error
	CreateAccount(ctx context.Context, account *model.Account) error
	Update(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, id uint) error
	GetAll(ctx context.Context, page, limit int) ([]model.User, int64, error)
	GetBySDT(ctx context.Context, sdt string) (*model.User, error)
	GetRoleByID(ctx context.Context, id uint) (*model.Role, error)
	ListRoles(ctx context.Context, page, limit int) ([]model.Role, int64, error)
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
func (r *UserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
	var user model.User
	result := r.GetDB().WithContext(ctx).Preload("Account.Role").First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// Tạo user
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	return generic.GenericCreate(ctx, r.GetDB(), user)
}

// Tạo account
func (r *UserRepository) CreateAccount(ctx context.Context, account *model.Account) error {
	return generic.GenericCreate(ctx, r.GetDB(), account)
}

// Cập nhật user
func (r *UserRepository) Update(ctx context.Context, user *model.User) error {
	return generic.GenericUpdate(ctx, r.GetDB(), user)
}

// Xóa user
func (r *UserRepository) Delete(ctx context.Context, id uint) error {
	return generic.GenericDelete[model.User, uint](ctx, r.GetDB(), id)
}

// Phân trang với offset/limit
func (r *UserRepository) GetAll(ctx context.Context, page, limit int) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	// Count total
	if err := r.GetDB().WithContext(ctx).Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	result := r.GetDB().WithContext(ctx).Preload("Account.Role").Offset(offset).Limit(limit).Find(&users)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return users, total, nil
}

// Tìm user theo số điện thoại với preload Account.Role
func (r *UserRepository) GetBySDT(ctx context.Context, sdt string) (*model.User, error) {
	var user model.User
	result := r.GetDB().WithContext(ctx).Preload("Account.Role").Where("sdt = ?", sdt).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) GetRoleByID(ctx context.Context, id uint) (*model.Role, error) {
	return generic.GenericGetByID[model.Role, uint](ctx, r.GetDB(), id)
}

func (r *UserRepository) ListRoles(ctx context.Context, page, limit int) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64
	if err := r.GetDB().WithContext(ctx).Model(&model.Role{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	result := r.GetDB().WithContext(ctx).Offset(offset).Limit(limit).Find(&roles)
	if result.Error != nil {
		return nil, 0, result.Error
	}
	return roles, total, nil
}
