package db

import (
	"time"

	"github.com/uptrace/bun"
)

type IndexTracker struct {
	bun.BaseModel `bun:"table:index_tracker"`

	Board        string    `bun:"board,pk"`
	LastModified time.Time `bun:"last_modified"`
	PostNumber   int64     `bun:"post_number"`
}
