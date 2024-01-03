package util

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config structure for reading YAML configuration
type Config struct {
	Postgres      Postgres      `yaml:"postgres"`
	ElasticSearch Elasticsearch `yaml:"elasticsearch"`
}

type Postgres struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
}

type Elasticsearch struct {
	URL string `yaml:"url"`
}

func ParseNewConfig(name string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
