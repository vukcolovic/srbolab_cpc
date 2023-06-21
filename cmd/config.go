package main

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	HTTP       string `yaml:"http"`
	HTTPPort   int    `yaml:"http_port"`
	DbHost     string `yaml:"db_host"`
	DbName     string `yaml:"db_name"`
	DbUser     string `yaml:"db_user"`
	DbPassword string `yaml:"db_password"`
}

func loadYamlConfig(conf *Config, fileName string) error {
	f, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer func() { _ = f.Close() }()
	return yaml.NewDecoder(f).Decode(conf)

}
