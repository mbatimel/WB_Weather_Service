package config

import ("os"
"gopkg.in/yaml.v3"
)

type Repo struct {
	Database string `yaml:"database"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
}

type ConfigWrapRepo struct {
	Config Repo `yaml:"repo"`
}

func NewConfigDB(path string) (*Repo, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := ConfigWrapRepo{}
	err = yaml.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config.Config, nil
}