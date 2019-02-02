package repositories

import (
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
	return nil
}
