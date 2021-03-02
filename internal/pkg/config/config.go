package config

import (
	"github.com/kelseyhightower/envconfig"
)

// LoadConfigs loads a configuration into an object
func LoadConfigs() (Config, error) {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Config set of configurations needed to run the app
type Config struct {
	Server struct {
		Port string `envconfig:"SERVER_PORT"`
	}
	Db struct {
		Host     string `envconfig:"DB_HOST"`
		Port     int    `envconfig:"DB_PORT"`
		User     string `envconfig:"DB_USER"`
		Password string `envconfig:"DB_PASSWORD"`
		Dbname   string `envconfig:"DB_NAME"`
		Driver   string `envconfig:"DB_DRIVER"`
	}
	App struct {
		KeySource struct {
			PoolSize   int `envconfig:"APP_KEYSOURCE_POOL_SIZE"`
			RSAKeySize int `envconfig:"APP_KEYSOURCE_RSAKEY_SIZE"`
		}
	}
}
