package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
)

type Config struct {
	ServerAddress string `yaml:"server_address"`

	Postgres struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		DBName   string `yaml:"db_name"`
		SSLMode  string `yaml:"ssl_mode"`
	} `yaml:"postgres"`
	Minio struct {
		Url                  string `yaml:"url"`
		User                 string `yaml:"user"`
		Password             string `yaml:"password"`
		Token                string `yaml:"token"`
		Ssl                  bool   `yaml:"ssl"`
		UserObjectBucketName string `yaml:"bucket_name"`
	} `yaml:"minio"`
}

func Read() (*Config, error) {
	var cfg Config
	path, err := filepath.Abs("config.yaml")
	if err != nil {
		log.Fatalf("creating path: %v", err)
		return nil, err
	}
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("can not load yaml: %v", err)
		return nil, err
	}
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("can not unmarshal config %v", err)
		return nil, err
	}
	return &cfg, nil
}
