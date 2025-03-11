package movies

import (
	"Movies-Go/internal/entity"
	basic_repo "Movies-Go/internal/repository/postgres/_basic_repo"
	"Movies-Go/internal/repository/postgres/movies"
	"context"
)

type Repository interface {
	Create(ctx context.Context, data movies.CreateMovieRequest) (entity.Movie, error)
	GetAll(ctx context.Context, filter movies.SearchMovieRequest) ([]*entity.Movie, int, error)
	GetByID(ctx context.Context, id int) (*entity.Movie, error)
	Update(ctx context.Context, data movies.UpdateMovieRequest) (entity.Movie, error)
	Delete(ctx context.Context, data basic_repo.Delete) error
	Search(ctx context.Context, query string, page, limit int) ([]*movies.MovieResponse, int, error)
}
