package config

type Config struct {
	Server ServerConfig `yaml:"server"`
	Repo   Repo   `yaml:"repo"`
}