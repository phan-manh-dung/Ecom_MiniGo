package repository

import (
	"gin/product_service/model"

	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetProduct(id uint) (*model.Product, error) {
	var product model.Product
	result := r.db.Find(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}
