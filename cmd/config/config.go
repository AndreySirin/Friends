package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	App App `yaml:"app"`
}
type Server struct {
	HTTPPort string `yaml:"httpPort"`
}
type Database struct {
	Address      string `yaml:"address"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	NameDatabase string `yaml:"nameDatabase"`
	Dialect      string `yaml:"dialect"`
	Datasource   string `yaml:"datasource"`
	Dir          string `yaml:"dir"`
	Table        string `yaml:"table"`
}
type Development struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
}
type App struct {
	Name        string      `yaml:"name"`
	Version     string      `yaml:"version"`
	Development Development `yaml:"development"`
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	fmt.Printf("Загруженная конфигурация: %+v\n", config)
	return &config, nil
}
