package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port           int    `envconfig:"PORT" default:"8080"`
	DatabaseURL    string `envconfig:"DATABASE_URL" required:"true"`
	LMSBaseURL     string `envconfig:"LMS_BASE_URL" required:"true"`
	PSPBaseURL     string `envconfig:"PSP_BASE_URL" required:"true"`
	ProductBaseURL string `envconfig:"PRODUCT_BASE_URL" required:"true"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
