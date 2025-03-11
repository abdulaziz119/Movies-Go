package users

import (
	"Movies-Go/internal/repository/postgres/movies"
	"Movies-Go/internal/repository/postgres/users"
	"Movies-Go/internal/util/password"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type Controller struct {
	repo Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	var filter movies.Filter
	if err := ctx.ShouldBindQuery(&filter); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid filter parameters: " + err.Error(),
		})
		return
	}

	users, err := c.repo.GetAll(ctx, filter)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": users,
	})
}

func (c *Controller) GetByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	user, err := c.repo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data": user,
	})
}

func (c *Controller) Update(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	existingUser, err := c.repo.GetByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "User not found",
		})
		return
	}

	var req users.UpdateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request data: " + err.Error(),
		})
		return
	}

	if req.Name != "" {
		existingUser.Name = req.Name
	}

	if req.Email != "" && req.Email != existingUser.Email {
		user, _ := c.repo.GetByEmail(ctx, req.Email)
		if user != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Email already in use",
			})
			return
		}
		existingUser.Email = req.Email
	}

	if req.Password != "" {
		hashedPassword, err := password.Hash(req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to process password: " + err.Error(),
			})
			return
		}
		existingUser.Password = hashedPassword
	}

	if err := c.repo.Update(ctx, existingUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User updated successfully",
		"data":    existingUser,
	})
}

func (c *Controller) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid user ID",
		})
		return
	}

	currentUserID, exists := ctx.Get("user_id")
	if exists && currentUserID.(int) == id {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Cannot delete your own account",
		})
		return
	}

	if err := c.repo.Delete(ctx, id); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete user: " + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}
