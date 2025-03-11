package movies

import (
	"Movies-Go/internal/controller/http/v1/movies"
	"Movies-Go/internal/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func Router(router *gin.RouterGroup, controller *movies.Controller) {
	moviesGroup := router.Group("/movies")

	moviesGroup.Use(middleware.AuthMiddleware())
	{
		moviesGroup.GET("", controller.GetAll)
		moviesGroup.GET("/:id", controller.GetByID)
		moviesGroup.GET("/search", controller.Search)
		moviesGroup.POST("", controller.Create)
		moviesGroup.PUT("/:id", controller.Update)
		moviesGroup.DELETE("/:id", controller.Delete)
	}
}
