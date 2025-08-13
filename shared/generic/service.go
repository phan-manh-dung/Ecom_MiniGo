package generic

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

// GenericService interface cơ bản cho tất cả service
type GenericService[T any, R any] interface {
	Execute(ctx context.Context, req T) (R, error)
}

// CRUDService interface cho các operation CRUD cơ bản
type CRUDService[T any, R any, ID any] interface {
	Get(ctx context.Context, id ID) (R, error)
	Create(ctx context.Context, req T) (R, error)
	Update(ctx context.Context, req T) (R, error)
	Delete(ctx context.Context, id ID) error
	List(ctx context.Context) ([]R, error)
}

// BaseService struct cơ bản cho tất cả service
type BaseService struct{}

// NewBaseService tạo instance mới của BaseService
func NewBaseService() *BaseService {
	return &BaseService{}
}

// ValidateRequest generic validation method
func ValidateRequest[T any](req T) error {
	// Implement validation logic here
	return nil
}

// LogOperation generic logging method
func (s *BaseService) LogOperation(operation string, details map[string]interface{}) {
	// Implement logging logic here
}

// Generic business logic functions
func HandleGetByID[T any, R any, ID any](
	ctx context.Context,
	id ID,
	getFunc func(ID) (*T, error),
	convertFunc func(*T) R,
	notFoundMessage string,
	successMessage string,
) (R, error) {
	entity, err := getFunc(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Return empty response for not found
			return *new(R), nil
		}
		return *new(R), fmt.Errorf("failed to get entity: %v", err)
	}

	// Convert to response
	response := convertFunc(entity)
	return response, nil
}

func HandleCreate[T any, R any](
	ctx context.Context,
	req T,
	createFunc func(*T) error,
	convertFunc func(*T) R,
	successMessage string,
) (R, error) {
	// Convert request to model
	entity := &req
	
	// Create entity
	if err := createFunc(entity); err != nil {
		return *new(R), fmt.Errorf("failed to create entity: %v", err)
	}

	// Convert to response
	response := convertFunc(entity)
	return response, nil
}

func HandleUpdate[T any, R any, ID any](
	ctx context.Context,
	id ID,
	req T,
	getFunc func(ID) (*T, error),
	updateFunc func(*T) error,
	convertFunc func(*T) R,
	successMessage string,
) (R, error) {
	// Get existing entity
	entity, err := getFunc(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return *new(R), fmt.Errorf("entity not found")
		}
		return *new(R), fmt.Errorf("failed to get entity: %v", err)
	}

	// Update entity fields (this needs to be implemented per entity type)
	// For now, just update the entity
	
	// Save updated entity
	if err := updateFunc(entity); err != nil {
		return *new(R), fmt.Errorf("failed to update entity: %v", err)
	}

	// Convert to response
	response := convertFunc(entity)
	return response, nil
}

func HandleDelete[T any, ID any](
	ctx context.Context,
	id ID,
	getFunc func(ID) (*T, error),
	deleteFunc func(ID) error,
	successMessage string,
) error {
	// Check if entity exists
	_, err := getFunc(id)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("entity not found")
		}
		return fmt.Errorf("failed to get entity: %v", err)
	}

	// Delete entity
	if err := deleteFunc(id); err != nil {
		return fmt.Errorf("failed to delete entity: %v", err)
	}

	return nil
}

func HandleList[T any, R any](
	ctx context.Context,
	req interface{},
	listFunc func(int, int) ([]T, int64, error),
	convertFunc func(*T) R,
	successMessage string,
) ([]R, int64, error) {
	// Default pagination
	page, limit := 1, 10
	
	// Try to extract page and limit from request
	// This is a simplified version, you might need to implement per request type
	
	// Get paginated results
	entities, total, err := listFunc(page, limit)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list entities: %v", err)
	}

	// Convert to response
	responses := make([]R, 0, len(entities))
	for _, entity := range entities {
		response := convertFunc(&entity)
		responses = append(responses, response)
	}

	return responses, total, nil
}
