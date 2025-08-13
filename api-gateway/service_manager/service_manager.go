package service_manager

import (
	"log"
)

// ServiceManager quản lý tất cả gRPC clients
type ServiceManager struct {
	UserClient    interface{} // Sẽ được inject từ Wire
	ProductClient interface{} // Sẽ được inject từ Wire
	OrderClient   interface{} // Sẽ được inject từ Wire
}

// NewServiceManager tạo instance mới của ServiceManager
func NewServiceManager() (*ServiceManager, error) {
	log.Printf("Initializing ServiceManager...")

	// Tạm thời return empty để Wire có thể inject dependencies
	return &ServiceManager{}, nil
}
