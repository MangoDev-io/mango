package config

import "github.com/kelseyhightower/envconfig"

// New returns an instance of Config
func New() (*Config, error) {
	var c Config
	err := envconfig.Process("API", &c)
	return &c, err
}
