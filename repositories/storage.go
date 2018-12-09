package repositories

import (
	"github.com/topfreegames/extensions/pg"
)

// Storage holds pointers to storage engines used by
// repositories
type Storage struct {
	PG *pg.Client
}

// NewStorage ctor
func NewStorage() *Storage {
	return &Storage{}
}
