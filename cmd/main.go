package main

import (
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"go.uber.org/fx"
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

func ProvideDB() *bun.DB {
	return postgres.NewPostgres()
}

func ProvideMoviesRepo(db *bun.DB) *movies.Repository {
	return movies.NewRepository(db)
}

func ProvideUsersRepo(db *bun.DB) *users.Repository {
	return users.NewRepository(db)
}

func ProvideMoviesController(repo *movies.Repository) *movies_controller.Controller {
	return movies_controller.NewController(repo)
}

func ProvideUsersController(repo *users.Repository) *users_controller.Controller {
	return users_controller.NewController(repo)
}

func ProvideAuthController(repo *users.Repository) *auth_controller.Controller {
	return auth_controller.NewController(repo)
}

func ProvideRouter() *gin.Engine {
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

	return r
}

func RegisterRoutes(
	r *gin.Engine,
	moviesController *movies_controller.Controller,
	usersController *users_controller.Controller,
	authController *auth_controller.Controller,
) {
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
}

// StartServer starts the HTTP server
func StartServer(lifecycle fx.Lifecycle, r *gin.Engine) {
	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Println("Starting server on port", config.GetConf().Port)
			go func() {
				if err := r.Run(":" + config.GetConf().Port); err != nil {
					log.Fatal("Server failed to start:", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping server")
			return nil
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			ProvideDB,
			ProvideMoviesRepo,
			ProvideUsersRepo,
			ProvideMoviesController,
			ProvideUsersController,
			ProvideAuthController,
			ProvideRouter,
		),
		fx.Invoke(RegisterRoutes, StartServer),
	).Run()
}
