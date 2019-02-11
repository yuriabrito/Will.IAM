package repositories

import (
	"strings"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
	"github.com/topfreegames/extensions/redis"
)

func loadDefaultConfigRedis(config *viper.Viper) {
	config.SetDefault("extensions.redis.url", "redis://localhost:6379")
	config.SetDefault("extensions.redis.connectionTimeout", 3)
}

// ConfigureRedis sets s.Redis
func (s *Storage) ConfigureRedis(config *viper.Viper) error {
	loadDefaultConfigRedis(config)
	client, err := redis.NewClient("extensions.redis", config, nil, nil)
	if err != nil {
		return err
	}
	s.Redis = client
	s.RedisPool = newRedisPool(config.GetString("extensions.redis.url"))
	return nil
}

func newRedisPool(url string) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:   2,
		MaxActive: 4,
		Dial: func() (redigo.Conn, error) {
			r, err := redigo.Dial("tcp", strings.Replace(url, "redis://", "", 1))
			if err != nil {
				println(err.Error())
			}
			return r, err
		},
	}
}
