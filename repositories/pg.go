package repositories

import (
	"github.com/spf13/viper"
	"github.com/topfreegames/extensions/pg"
)

func loadDefaultConfigPG(config *viper.Viper) {
	config.SetDefault("extensions.pg.user", "postgres")
	config.SetDefault("extensions.pg.pass", "Will.IAM")
	config.SetDefault("extensions.pg.host", "localhost")
	config.SetDefault("extensions.pg.database", "Will.IAM")
	config.SetDefault("extensions.pg.port", 8432)
	config.SetDefault("extensions.pg.poolSize", 20)
	config.SetDefault("extensions.pg.maxRetries", 3)
	config.SetDefault("extensions.pg.connectionTimeout", 5)
}

// ConfigurePG sets s.PG
func (s *Storage) ConfigurePG(config *viper.Viper) error {
	loadDefaultConfigPG(config)
	client, err := pg.NewClient("extensions.pg", config, nil, nil)
	if err != nil {
		return err
	}
	s.PG = client
	return nil
}
