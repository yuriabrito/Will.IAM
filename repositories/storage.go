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

// Clone storage
func (s *Storage) Clone() *Storage {
	var pg *pg.Client
	*pg = *s.PG
	// TODO: copy DB interface
	return &Storage{PG: pg}
}

type withStorage struct {
	storage *Storage
}

func (ws *withStorage) setStorage(s *Storage) {
	ws.storage = s
}
