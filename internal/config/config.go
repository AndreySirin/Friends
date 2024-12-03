package config

import (
	"fmt"
	"log/slog"
	"os"

	"gopkg.in/yaml.v3"
)

// TODO: добавим валидацию на конфиг

type Config struct {
	App App `yaml:"app"`
}

type App struct {
	Name        string      `yaml:"name"`
	Version     string      `yaml:"version"`
	Development Development `yaml:"development"`
}

type Development struct {
	Server   Server   `yaml:"server"`
	Database Database `yaml:"database"`
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

func LoadConfig(lg *slog.Logger, path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	var config Config
	if err = yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	lg.With("module", "config").Debug(fmt.Sprintf("%+v\n", config))

	return &config, nil
}
