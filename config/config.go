package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type ApiConfig struct {
	Database DatabaseConfig `json:"database"`
	ApiHost  string         `json:"api_host"`
	ApiPort  string         `json:"api_port"`
}

const (
	configVar = "API_CONFIG"
)

var config *ApiConfig

// LoadConfig tries to load the project configuration
func LoadConfig() error {
	configPath := "/etc/go-api/config.json"

	if path, ok := os.LookupEnv(configVar); ok && path != "" {
		configPath = path
	} else {
		log.Println("`API_CONFIG` not set or empty, using default path: ", configPath)
	}

	raw, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Println(err)
		return err
	}

	if err := json.Unmarshal(raw, &config); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetConfig returns the config data
func GetConfig() *ApiConfig {
	if config == nil {
		log.Println("Configuration not loaded")
	}

	return config
}
