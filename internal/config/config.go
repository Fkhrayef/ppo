package config

import "github.com/kelseyhightower/envconfig"

type Config struct {
	Port           int    `envconfig:"PORT" default:"8080"`
	DatabaseURL    string `envconfig:"DATABASE_URL" required:"true"`
	UseFakeClients bool   `envconfig:"USE_FAKE_CLIENTS" default:"true"`
	LMSBaseURL     string `envconfig:"LMS_BASE_URL" default:"http://localhost:8081"`
	PSPBaseURL     string `envconfig:"PSP_BASE_URL" default:"http://localhost:8082"`
	ProductBaseURL string `envconfig:"PRODUCT_BASE_URL" default:"http://localhost:8083"`
}

func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
