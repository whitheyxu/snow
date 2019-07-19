// configj
package config

type Config struct {
	IsDebugAccess bool
}

var Cfg *Config

const (
	defaultIsDebugAccess = false
)

func init() {
	Cfg = new(Config)
	Cfg.IsDebugAccess = defaultIsDebugAccess
}

func GetConfig() *Config {
	return Cfg
}

func SetIsDebugAccess(isDebugAccess bool) {
	Cfg.IsDebugAccess = isDebugAccess
}
