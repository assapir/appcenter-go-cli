package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// Config - The representation of the config
type Config struct {
	OwnerName string `json:"owner"`
	AppName   string `json:"app"`
	APIKey    string `json:"apiKey"`
}

// GetConfig - Deserialize the config from a 'config.json' file
func GetConfig() (*Config, error) {
	pwd, _ := os.Getwd()
	data, err := ioutil.ReadFile(pwd + "/config.json")
	if err != nil {
		return nil, err
	}

	var jsonConfig Config
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return nil, err
	}

	return &jsonConfig, nil
}
