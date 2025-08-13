//go:build wireinject
// +build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

// wireApp khởi tạo toàn bộ application với dependency injection
func wireApp() (*gin.Engine, error) {
	wire.Build(
		// tạo gin engine đơn giản
		provideGinEngine,
	)
	return &gin.Engine{}, nil
}

// provideGinEngine tạo gin engine đơn giản
func provideGinEngine() *gin.Engine {
	return gin.New()
}
