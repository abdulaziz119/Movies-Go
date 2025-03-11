package users

import (
	"Movies-Go/internal/entity"
	"Movies-Go/internal/repository/postgres/movies"
	"context"
)

type Repository interface {
	Login(ctx context.Context, email, password string) (*entity.User, error)

	GetByEmail(ctx context.Context, email string) (*entity.User, error)

	Create(ctx context.Context, user *entity.User) error

	GetAll(ctx context.Context, filter movies.Filter) ([]*entity.User, error)

	GetByID(ctx context.Context, id int) (*entity.User, error)

	Update(ctx context.Context, user *entity.User) error

	Delete(ctx context.Context, id int) error
}
