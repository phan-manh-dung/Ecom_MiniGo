package repository

import (
	"gin/product_service/model"
	"gin/shared/generic"

	"gorm.io/gorm"
)

type ProductRepository struct {
	*generic.BaseRepository
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{
		BaseRepository: generic.NewBaseRepository(db),
	}
}

// GetProduct - alias cho GetByID để tương thích với service
func (r *ProductRepository) GetProduct(id uint) (*model.Product, error) {
	return generic.GenericGetByID[model.Product, uint](r.GetDB(), id)
}

func (r *ProductRepository) GetByID(id uint) (*model.Product, error) {
	var product model.Product
	result := r.GetDB().First(&product, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &product, nil
}

func (r *ProductRepository) Create(product *model.Product) error {
	return generic.GenericCreate(r.GetDB(), product)
}

func (r *ProductRepository) Update(product *model.Product) error {
	return generic.GenericUpdate(r.GetDB(), product)
}

func (r *ProductRepository) Delete(id uint) error {
	return generic.GenericDelete[model.Product, uint](r.GetDB(), id)
}

func (r *ProductRepository) DecreaseInventory(productId uint, quantity int) error {
	var inventory model.Inventory
	if err := r.GetDB().Where("product_id = ?", productId).First(&inventory).Error; err != nil {
		return err
	}
	return r.GetDB().Model(&inventory).Update("quantity", inventory.Quantity-quantity).Error
}

func (r *ProductRepository) IncreaseInventory(productId uint, quantity int) error {
	var inventory model.Inventory
	if err := r.GetDB().Where("product_id = ?", productId).First(&inventory).Error; err != nil {
		return err
	}
	return r.GetDB().Model(&inventory).Update("quantity", inventory.Quantity+quantity).Error
}
