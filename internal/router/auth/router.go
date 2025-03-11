package auth

import (
	"Movies-Go/internal/controller/http/v1/auth"
	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup, controller *auth.Controller) {
	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", controller.Register)
		authGroup.POST("/login", controller.Login)
	}
}
