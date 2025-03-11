package users

import (
	"Movies-Go/internal/entity"
	"Movies-Go/internal/repository/postgres/movies"
	"context"
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

func (r *Repository) Create(ctx context.Context, user *entity.User) error {
	now := time.Now()
	user.CreatedAt = &now
	user.UpdatedAt = &now

	_, err := r.db.NewInsert().Model(user).Exec(ctx)
	return err
}

func (r *Repository) GetAll(ctx context.Context, filter movies.Filter) ([]*entity.User, error) {
	var users []*entity.User

	query := r.db.NewSelect().
		Model(&users).
		Where("deleted_at IS NULL").
		Order("id ASC")

	if filter.Page != nil && filter.Limit != nil {
		page := *filter.Page
		limit := *filter.Limit

		if page < 1 {
			page = 1
		}

		if limit < 1 {
			limit = 10
		}

		offset := (page - 1) * limit
		query = query.Limit(limit).Offset(offset)
	}

	err := query.Scan(ctx)
	return users, err
}

func (r *Repository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	user := new(entity.User)

	err := r.db.NewSelect().
		Model(user).
		Where("id = ? AND deleted_at IS NULL", id).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	user := new(entity.User)

	err := r.db.NewSelect().
		Model(user).
		Where("email = ? AND deleted_at IS NULL", email).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *Repository) Update(ctx context.Context, user *entity.User) error {
	now := time.Now()
	user.UpdatedAt = &now

	_, err := r.db.NewUpdate().
		Model(user).
		Where("id = ? AND deleted_at IS NULL", user.Id).
		Exec(ctx)

	return err
}

func (r *Repository) Delete(ctx context.Context, id int) error {
	now := time.Now()

	_, err := r.db.NewUpdate().
		Model((*entity.User)(nil)).
		Set("deleted_at = ?", now).
		Where("id = ? AND deleted_at IS NULL", id).
		Exec(ctx)

	return err
}

func (r *Repository) Login(ctx context.Context, email, password string) (*entity.User, error) {
	return r.GetByEmail(ctx, email)
}

func (r *Repository) Register(ctx context.Context, user *entity.User) error {
	return r.Create(ctx, user)
}
