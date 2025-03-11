package entity

import (
	"github.com/uptrace/bun"
	"time"
)

type Movie struct {
	bun.BaseModel `bun:"table:movies"`

	basicEntity
	Id        int        `json:"id" bun:"id,pk,autoincrement"`
	Title     string     `json:"title" bun:"title,notnull"`
	Director  string     `json:"director" bun:"director,notnull"`
	Year      int        `json:"year" bun:"year,notnull"`
	Plot      string     `json:"plot" bun:"plot"`
	Rating    float64    `json:"rating" bun:"rating"`
	CreatedAt *time.Time `json:"created_at" bun:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bun:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" bun:"deleted_at"`
}
