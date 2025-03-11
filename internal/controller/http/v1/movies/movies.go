package movies

import (
	basic_controller "Movies-Go/internal/controller/http/v1/_basic_controller"
	"Movies-Go/internal/entity"
	basic_repo "Movies-Go/internal/repository/postgres/_basic_repo"
	"Movies-Go/internal/repository/postgres/movies"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type MovieRepositoryAdapter struct {
	repo *movies.Repository
}

func (a *MovieRepositoryAdapter) Create(ctx context.Context, data movies.CreateMovieRequest) (entity.Movie, error) {
	movie := &entity.Movie{}

	if data.Title != nil {
		movie.Title = *data.Title
	}

	if data.Director != nil {
		movie.Director = *data.Director
	}

	if data.Year != nil {
		movie.Year = *data.Year
	}

	if data.Plot != nil {
		movie.Plot = *data.Plot
	}

	if data.Rating != nil {
		movie.Rating = *data.Rating
	}

	err := a.repo.Create(ctx, movie)
	if err != nil {
		return entity.Movie{}, err
	}

	return *movie, nil
}

func (a *MovieRepositoryAdapter) GetAll(ctx context.Context, filter movies.SearchMovieRequest) ([]*entity.Movie, int, error) {
	return a.repo.GetAll(ctx, filter)
}

func (a *MovieRepositoryAdapter) GetByID(ctx context.Context, id int) (*entity.Movie, error) {
	return a.repo.GetByID(ctx, id)
}

func (a *MovieRepositoryAdapter) Update(ctx context.Context, data movies.UpdateMovieRequest) (entity.Movie, error) {
	id := 0
	if data.Id != nil {
		id = *data.Id
	}

	movie, err := a.repo.GetByID(ctx, id)
	if err != nil {
		return entity.Movie{}, err
	}

	if data.Title != nil {
		movie.Title = *data.Title
	}

	if data.Director != nil {
		movie.Director = *data.Director
	}

	if data.Year != nil {
		movie.Year = *data.Year
	}

	if data.Plot != nil {
		movie.Plot = *data.Plot
	}

	if data.Rating != nil {
		movie.Rating = *data.Rating
	}

	err = a.repo.Update(ctx, movie)
	if err != nil {
		return entity.Movie{}, err
	}

	return *movie, nil
}

func (a *MovieRepositoryAdapter) Delete(ctx context.Context, data basic_repo.Delete) error {
	if data.Id == nil {
		return fmt.Errorf("id is required")
	}
	return a.repo.Delete(ctx, data)
}

func (a *MovieRepositoryAdapter) Search(ctx context.Context, query string, page, limit int) ([]*movies.MovieResponse, int, error) {
	filter := movies.SearchMovieRequest{
		Query: &query,
		Page:  &page,
		Limit: &limit,
	}

	moviesResult, count, err := a.repo.GetAll(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	var response []*movies.MovieResponse
	for _, movie := range moviesResult {
		response = append(response, &movies.MovieResponse{
			ID:        &movie.Id,
			Title:     &movie.Title,
			Director:  &movie.Director,
			Year:      &movie.Year,
			Plot:      &movie.Plot,
			Rating:    &movie.Rating,
			CreatedAt: movie.CreatedAt,
			UpdatedAt: movie.UpdatedAt,
		})
	}

	return response, count, nil
}

type Controller struct {
	useCase Repository
}

func NewController(repo *movies.Repository) *Controller {
	adapter := &MovieRepositoryAdapter{
		repo: repo,
	}
	return &Controller{
		useCase: adapter,
	}
}

func (cl *Controller) Create(c *gin.Context) {
	var request movies.CreateMovieRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	detail, err := cl.useCase.Create(c.Request.Context(), request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": detail,
	})
}

func (cl *Controller) GetAll(c *gin.Context) {
	var filter movies.SearchMovieRequest
	query := c.Request.URL.Query()

	limitQ := query["limit"]
	if len(limitQ) > 0 {
		queryInt, err := strconv.Atoi(limitQ[0])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"message": "Limit must be number!",
				"status":  false,
			})

			return
		}

		filter.Limit = &queryInt
	}

	pageQ := query["page"]
	if len(pageQ) > 0 {
		page, err := strconv.Atoi(pageQ[0])
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Page must be number!",
				"status":  false,
			})
			return
		}
		filter.Page = &page
	}

	queryQ := query["query"]
	if len(queryQ) > 0 {
		filter.Query = &queryQ[0]
	}

	ctx := context.Background()

	list, count, err := cl.useCase.GetAll(ctx, filter)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data": map[string]interface{}{
			"results": list,
			"count":   count,
		},
	})
}

func (cl *Controller) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid movie ID",
			"status": false,
		})
		return
	}
	ctx := context.Background()

	detail, err := cl.useCase.GetByID(ctx, id)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl *Controller) Update(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":  "Invalid movie ID",
			"status": false,
		})
		return
	}

	var data movies.UpdateMovieRequest
	err = c.ShouldBind(&data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	if data.Id == nil {
		data.Id = &id
	}
	ctx := context.Background()

	detail, err := cl.useCase.Update(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
		"data":    detail,
	})
}

func (cl *Controller) Delete(c *gin.Context) {
	ctx, data, err := basic_controller.BasicDelete(c)
	if err != nil {
		return
	}

	err = cl.useCase.Delete(ctx, data)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": err.Error(),
			"status":  false,
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ok!",
		"status":  true,
	})
}

// Search handles movie search requests
func (cl *Controller) Search(c *gin.Context) {
	var request movies.SearchMovieRequest
	if err := c.ShouldBindQuery(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := ""
	if request.Query != nil {
		query = *request.Query
	}

	page := 1
	if request.Page != nil {
		page = *request.Page
	}

	limit := 10
	if request.Limit != nil {
		limit = *request.Limit
	}

	moviesResult, totalCount, err := cl.useCase.Search(c.Request.Context(), query, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := (totalCount + limit - 1) / limit
	if totalPages < 1 {
		totalPages = 1
	}

	response := gin.H{
		"data":        moviesResult,
		"total":       totalCount,
		"page":        page,
		"limit":       limit,
		"total_pages": totalPages,
	}

	c.JSON(http.StatusOK, response)
}
