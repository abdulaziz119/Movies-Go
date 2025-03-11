package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	basicEntity
	Id        int        `json:"id" bun:"id,pk,autoincrement"`
	Name      string     `json:"name" bun:"name,notnull"`
	Email     string     `json:"email" bun:"email,notnull,unique"`
	Password  string     `json:"-" bun:"password,notnull"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bun:"deleted_at"`
}
