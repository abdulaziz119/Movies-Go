package auth

import (
	"Movies-Go/internal/entity"
	"Movies-Go/internal/pkg/auth"
	"Movies-Go/internal/repository/postgres/users"
	"Movies-Go/internal/util/password"
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Repository interface {
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) error
}

type Controller struct {
	userRepo Repository
}

func NewController(userRepo Repository) *Controller {
	return &Controller{
		userRepo: userRepo,
	}
}

func (c *Controller) Register(ctx *gin.Context) {
	var req users.RegisterRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	existingUser, err := c.userRepo.GetByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "User with this email already exists",
		})
		return
	}

	hashedPassword, err := password.Hash(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to process password: " + err.Error(),
		})
		return
	}

	now := time.Now()
	user := &entity.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	if err := c.userRepo.Create(ctx, user); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	token, err := auth.GenerateToken(user.Id, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token: " + err.Error(),
		})
		return
	}

	response := users.AuthResponse{
		Token: token,
		User: users.UserResponse{
			ID:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) Login(ctx *gin.Context) {
	var req users.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	user, err := c.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if !password.Verify(user.Password, req.Password) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token, err := auth.GenerateToken(user.Id, user.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token: " + err.Error(),
		})
		return
	}

	response := users.AuthResponse{
		Token: token,
		User: users.UserResponse{
			ID:        user.Id,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}

	ctx.JSON(http.StatusOK, response)
}
