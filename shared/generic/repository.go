package generic

import (
	"context"

	"gorm.io/gorm"
)

// GenericRepository interface cơ bản cho tất cả repository
type GenericRepository[T any, ID any] interface {
	GetByID(ctx context.Context, id ID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]*T, error)
}

// CRUDRepository interface cho các operation CRUD cơ bản
type CRUDRepository[T any, ID any] interface {
	GetByID(ctx context.Context, id ID) (*T, error)
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id ID) error
	GetAll(ctx context.Context, page, limit int) ([]T, int64, error)
}

// BaseRepository struct cơ bản cho tất cả repository
type BaseRepository struct {
	db *gorm.DB
}

// NewBaseRepository tạo instance mới của BaseRepository
func NewBaseRepository(db *gorm.DB) *BaseRepository {
	return &BaseRepository{db: db}
}

// ValidateID generic ID validation
func ValidateID[ID any](id ID) error {
	// Implement ID validation logic here
	return nil
}

// LogRepositoryOperation generic logging method
func (r *BaseRepository) LogRepositoryOperation(ctx context.Context, operation string, details map[string]interface{}) {
	// Implement logging logic here
	// Có thể lấy request ID, user ID từ context để log
}

// GetDB trả về database instance
func (r *BaseRepository) GetDB() *gorm.DB {
	return r.db
}

// Generic CRUD operations - sử dụng function thay vì method
func GenericGetByID[T any, ID any](ctx context.Context, db *gorm.DB, id ID) (*T, error) {
	var entity T
	result := db.WithContext(ctx).First(&entity, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}

func GenericCreate[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Create(entity).Error
}

func GenericUpdate[T any](ctx context.Context, db *gorm.DB, entity *T) error {
	return db.WithContext(ctx).Save(entity).Error
}

func GenericDelete[T any, ID any](ctx context.Context, db *gorm.DB, id ID) error {
	var entity T
	return db.WithContext(ctx).Delete(&entity, id).Error
}

func GenericGetAll[T any](ctx context.Context, db *gorm.DB, page, limit int) ([]T, int64, error) {
	var entities []T
	var total int64

	// Count total
	if err := db.WithContext(ctx).Model(new(T)).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * limit
	result := db.WithContext(ctx).Offset(offset).Limit(limit).Find(&entities)
	if result.Error != nil {
		return nil, 0, result.Error
	}

	return entities, total, nil
}

func GenericGetByField[T any](ctx context.Context, db *gorm.DB, field string, value interface{}) (*T, error) {
	var entity T
	result := db.WithContext(ctx).Where(field+" = ?", value).First(&entity)
	if result.Error != nil {
		return nil, result.Error
	}
	return &entity, nil
}
