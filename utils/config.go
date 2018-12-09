package utils

import (
	"strings"

	"github.com/spf13/viper"
)

func getConfig(path, envPrefix, configType string) (*viper.Viper, error) {
	config := viper.New()
	config.SetConfigFile(path)
	config.SetConfigType(configType)
	config.SetEnvPrefix(envPrefix)
	config.AddConfigPath(".")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	config.AutomaticEnv()

	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// GetConfig reads config from and returns it as viper config
func GetConfig(path string) (*viper.Viper, error) {
	config, err := getConfig(path, "Will.IAM", "yaml")
	if err != nil {
		return nil, err
	}
	return config, nil
}
