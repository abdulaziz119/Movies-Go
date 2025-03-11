package users

import (
	"Movies-Go/internal/controller/http/v1/users"
	"Movies-Go/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup, controller *users.Controller) {
	usersGroup := router.Group("/users")
	{
		usersGroup.Use(middleware.AuthMiddleware())
		{
			usersGroup.GET("", controller.GetAll)
			usersGroup.GET("/:id", controller.GetByID)

			adminGroup := usersGroup.Group("")
			{
				adminGroup.PUT("/:id", controller.Update)
				adminGroup.DELETE("/:id", controller.Delete)
			}
		}
	}
}
