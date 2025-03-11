package movies

import "time"

type CreateMovieRequest struct {
	Title    *string  `json:"title" binding:"required"`
	Director *string  `json:"director" binding:"required"`
	Year     *int     `json:"year" binding:"required,min=1800,max=2100"`
	Plot     *string  `json:"plot"`
	Rating   *float64 `json:"rating" binding:"min=0,max=10"`
}

type UpdateMovieRequest struct {
	Id       *int     `json:"id" form:"id"`
	Title    *string  `json:"title"`
	Director *string  `json:"director"`
	Year     *int     `json:"year" binding:"omitempty,min=1800,max=2100"`
	Plot     *string  `json:"plot"`
	Rating   *float64 `json:"rating" binding:"omitempty,min=0,max=10"`
}

type MovieResponse struct {
	ID        *int       `json:"id"`
	Title     *string    `json:"title"`
	Director  *string    `json:"director"`
	Year      *int       `json:"year"`
	Plot      *string    `json:"plot,omitempty"`
	Rating    *float64   `json:"rating"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type Filter struct {
	Page  *int `form:"page,default=1" binding:"min=1"`
	Limit *int `form:"limit,default=10" binding:"min=1,max=100"`
}

type SearchMovieRequest struct {
	Query *string `form:"query" binding:"required"`
	Page  *int    `form:"page,default=1" binding:"min=1"`
	Limit *int    `form:"limit,default=10" binding:"min=1,max=100"`
}

type SearchMovieResponse struct {
	Movies     []*MovieResponse `json:"data"`
	TotalCount int              `json:"total"`
	Page       int              `json:"page"`
	Limit      int              `json:"limit"`
	TotalPages int              `json:"total_pages"`
}
