package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"

	auth_controller "Movies-Go/internal/controller/http/v1/auth"
	movies_controller "Movies-Go/internal/controller/http/v1/movies"
	users_controller "Movies-Go/internal/controller/http/v1/users"
	"Movies-Go/internal/pkg/config"
	"Movies-Go/internal/pkg/repository/postgres"
	"Movies-Go/internal/repository/postgres/movies"
	"Movies-Go/internal/repository/postgres/users"
	auth_router "Movies-Go/internal/router/auth"
	movies_router "Movies-Go/internal/router/movies"
	users_router "Movies-Go/internal/router/users"
)

func main() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))

	postgresDB := postgres.NewPostgres()

	moviesRepo := movies.NewRepository(postgresDB)
	usersRepo := users.NewRepository(postgresDB)

	moviesController := movies_controller.NewController(moviesRepo)
	usersController := users_controller.NewController(usersRepo)
	authController := auth_controller.NewController(usersRepo)

	api := r.Group("api")
	{
		v1 := api.Group("v1")

		v1.GET("/health", func(c *gin.Context) {
			now := time.Now()
			c.JSON(http.StatusOK, gin.H{
				"status":  "healthy",
				"time":    now.Format(time.RFC3339),
				"version": "1.0.0",
			})
		})

		movies_router.Router(v1, moviesController)
		users_router.Router(v1, usersController)
		auth_router.Router(v1, authController)
	}

	log.Println("Starting server on port", config.GetConf().Port)
	log.Fatalln(r.Run(":" + config.GetConf().Port))
}
