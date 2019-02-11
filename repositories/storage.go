package repositories

import (
	redigo "github.com/gomodule/redigo/redis"
	"github.com/topfreegames/extensions/pg"
	"github.com/topfreegames/extensions/redis"
)

// Storage holds pointers to storage engines used by
// repositories
type Storage struct {
	PG        *pg.Client
	Redis     *redis.Client
	RedisPool *redigo.Pool
}

// NewStorage ctor
func NewStorage() *Storage {
	return &Storage{}
}

// Clone storage
func (s *Storage) Clone() *Storage {
	pg := &pg.Client{}
	redis := &redis.Client{}
	redisPool := &redigo.Pool{}
	if s.PG != nil {
		*pg = *s.PG
	}
	if s.Redis != nil {
		*redis = *s.Redis
	}
	if s.RedisPool != nil {
		*redisPool = *s.RedisPool
	}
	return &Storage{PG: pg, Redis: redis, RedisPool: redisPool}
}

type withStorage struct {
	storage *Storage
}

func (ws *withStorage) setStorage(s *Storage) {
	ws.storage = s
}
