package config

import (
	"fmt"
	"github.com/caarlos0/env/v9"
	"net"
	"net/url"
)

type Config struct {
	Host string `env:"HOST"`
	Port string `env:"PORT"`
}

func (c *Config) ServerAddress(endpoint string) string {
	u := url.URL{
		Scheme: "ws",
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   endpoint,
	}
	return u.String()
}

func Read() (*Config, error) {
	var config Config

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}
	return &config, nil
}
