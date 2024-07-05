package config

import (
	"os"
	"gopkg.in/yaml.v3"
)

type ServerConfig struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}
type ConfigWrapServer struct {
	Config ServerConfig `yaml:"server"`
}
func NewConfigsServer(path string) (*ServerConfig, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := ConfigWrapServer{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config.Config, nil
}