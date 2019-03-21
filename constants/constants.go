package constants

import "github.com/spf13/viper"

// Metrics constants
var Metrics = struct {
	ResponseTime string
}{
	ResponseTime: "response_time",
}

// AppInfo constants
var AppInfo = struct {
	Name    string
	Version string
}{
	Name:    "Will.IAM",
	Version: "1.0",
}

// constants from config
var (
	TokensCacheTTL             int
	TokensCacheEnabled         bool
	DefaultListOptionsPageSize int
)

// Set is called at start.Run
func Set(config *viper.Viper) {
	TokensCacheTTL = config.GetInt("tokens.cacheTTL")
	TokensCacheEnabled = config.GetBool("tokens.enabled")
	DefaultListOptionsPageSize = config.GetInt("listOptions.defaultPageSize")
}
