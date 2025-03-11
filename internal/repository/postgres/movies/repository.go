package movies

import (
	"Movies-Go/internal/entity"
	basic_repo "Movies-Go/internal/repository/postgres/_basic_repo"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/uptrace/bun"
)

type Repository struct {
	db *bun.DB
}

func NewRepository(db *bun.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(ctx context.Context, movie *entity.Movie) error {
	now := time.Now()
	movie.CreatedAt = &now
	movie.UpdatedAt = &now

	_, err := r.db.NewInsert().Model(movie).Exec(ctx)
	return err
}

func (r *Repository) GetByID(ctx context.Context, id int) (*entity.Movie, error) {
	movie := new(entity.Movie)

	err := r.db.NewSelect().
		Model(movie).
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *Repository) Update(ctx context.Context, movie *entity.Movie) error {
	now := time.Now()
	movie.UpdatedAt = &now

	_, err := r.db.NewUpdate().
		Model(movie).
		Where("id = ? AND deleted_at IS NULL", movie.Id).
		Exec(ctx)

	return err
}

func (r Repository) Delete(ctx context.Context, data basic_repo.Delete) error {
	return basic_repo.BasicDelete(ctx, data, &entity.User{}, r.db)
}

func (r *Repository) GetAll(ctx context.Context, filter SearchMovieRequest) ([]*entity.Movie, int, error) {
	page := 1
	if filter.Page != nil && *filter.Page > 0 {
		page = *filter.Page
	}

	limit := 10
	if filter.Limit != nil && *filter.Limit > 0 {
		limit = *filter.Limit
	}

	offset := (page - 1) * limit

	var movies []*entity.Movie
	var count int

	var query string
	if filter.Query != nil {
		query = *filter.Query
	}

	searchTerms := strings.Split(query, " ")
	var conditions []string
	var params []interface{}

	for _, term := range searchTerms {
		if term == "" {
			continue
		}
		term = "%" + strings.ToLower(term) + "%"
		conditions = append(conditions, "LOWER(title) LIKE ? OR LOWER(director) LIKE ? OR LOWER(plot) LIKE ?")
		params = append(params, term, term, term)
	}

	var whereClause string
	if len(conditions) > 0 {
		whereClause = "(" + strings.Join(conditions, " OR ") + ") AND deleted_at IS NULL"
	} else {
		whereClause = "deleted_at IS NULL"
	}

	countQuery := r.db.NewSelect().Model((*entity.Movie)(nil)).Where(whereClause, params...)
	count, err := countQuery.Count(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("error counting search results: %w", err)
	}

	err = r.db.NewSelect().
		Model(&movies).
		Where(whereClause, params...).
		Order("id ASC").
		Limit(limit).
		Offset(offset).
		Scan(ctx)

	if err != nil {
		return nil, 0, fmt.Errorf("error searching movies: %w", err)
	}

	return movies, count, nil
}
