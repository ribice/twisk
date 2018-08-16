package config

import (
	"fmt"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

// Load loads the configuration file from the given path
func Load(path string) (*Configuration, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file, %s", err)
	}

	var cfg Configuration
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct, %v", err)
	}

	return &cfg, nil
}

// Configuration holds application configuration data
type Configuration struct {
	Server  Server      `yaml:"server,omitempty"`
	DB      Database    `yaml:"database,omitempty"`
	JWT     JWT         `yaml:"jwt,omitempty"`
	App     Application `yaml:"application,omitempty"`
	OpenAPI OpenAPI     `yaml:"openapi,omitempty"`
}

// Database holds data necessery for database configuration
type Database struct {
	PSN            string `yaml:"psn,omitempty"`
	LogQueries     bool   `yaml:"log_queries,omitempty"`
	TimeoutSeconds int    `yaml:"timeout_seconds,omitempty"`
}

// Server holds data necessery for server configuration
type Server struct {
	Port                string `yaml:"port,omitempty"`
	ReadTimeoutSeconds  int    `yaml:"read_timeout_seconds,omitempty"`
	WriteTimeoutSeconds int    `yaml:"write_timeout_seconds,omitempty"`
}

// JWT holds data necessery for JWT configuration
type JWT struct {
	Secret    string `yaml:"secret,omitempty"`
	Duration  int    `yaml:"duration_minutes,omitempty"`
	Algorithm string `yaml:"signing_algorithm,omitempty"`
}

// Application represents application specific configuration
type Application struct {
	MinPasswordStrength int `yaml:"min_password_strength,omitempty"`
}

// OpenAPI holds username password for viewing api docs
type OpenAPI struct {
	Username string `yaml:"username,omitempty"`
	Password string `yaml:"password,omitempty"`
}
