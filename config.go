package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

// -------------
// CONFIGURATION
// -------------

// Config struct to hold the aliases from the config file
type Config struct {
	Aliases map[string]string `json:"aliases"`
}

// getConfigPath retrieves the path to the configuration file from the environment variable "AWOL_CONFIG_PATH".
// If the environment variable is not set, it defaults to "~/.config/awol/config.json".
func getConfigPath() string {
	cfgPath := os.Getenv("AWOL_CONFIG_PATH")
	if cfgPath == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "~/.config/awol/config.json"
		}
		cfgPath = path.Join(home, ".config/awol/config.json")
	}
	return cfgPath
}

// Reads the configuration file from the specified path and returns the Config struct
func readConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the config file does not exist, return an empty config without error
			return Config{Aliases: map[string]string{}}, nil
		}
		fmt.Fprintf(os.Stderr, "Error opening config file: %v\n", err)
		return Config{}, err
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error decoding config file: %v\n", err)
		return Config{}, err
	}

	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}

	return config, nil
}
