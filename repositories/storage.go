package repositories

import (
	"github.com/topfreegames/extensions/pg"
	"github.com/topfreegames/extensions/redis"
)

// Storage holds pointers to storage engines used by
// repositories
type Storage struct {
	PG    *pg.Client
	Redis *redis.Client
}

// NewStorage ctor
func NewStorage() *Storage {
	return &Storage{}
}

// Clone storage
func (s *Storage) Clone() *Storage {
	pg := &pg.Client{}
	redis := &redis.Client{}
	if s.PG != nil {
		*pg = *s.PG
	}
	if s.Redis != nil {
		*redis = *s.Redis
	}
	return &Storage{PG: pg, Redis: redis}
}

type withStorage struct {
	storage *Storage
}

func (ws *withStorage) setStorage(s *Storage) {
	ws.storage = s
}
