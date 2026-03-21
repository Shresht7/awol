package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
)

// -------------
// CONFIGURATION
// -------------

// Config struct to hold the aliases from the config file
type Config struct {
	BroadcastAddress string            `json:"broadcast"`
	Port             int               `json:"port"`
	Aliases          map[string]string `json:"aliases"`
}

func defaultConfig() Config {
	return Config{
		BroadcastAddress: "255.255.255.255",
		Port:             9,
		Aliases:          make(map[string]string),
	}
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
		cfgPath = path.Join(home, ".config", "awol", "config.json")
	}
	return cfgPath
}

// Reads the configuration file from the specified path and returns the Config struct
func readConfig(cfgPath string) (Config, error) {
	file, err := os.Open(cfgPath)
	if err != nil {
		if os.IsNotExist(err) {
			// If the config file does not exist, return an empty config without error
			return defaultConfig(), nil
		}
		return Config{}, fmt.Errorf("Error opening config file: %w", err)
	}
	defer file.Close()

	var config Config
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&config); err != nil {
		return Config{}, fmt.Errorf("Error decoding config file: %w", err)
	}

	defaultCfg := defaultConfig()
	if config.BroadcastAddress == "" {
		config.BroadcastAddress = defaultCfg.BroadcastAddress
	}
	if config.Port == 0 {
		config.Port = defaultCfg.Port
	}
	if config.Aliases == nil {
		config.Aliases = defaultCfg.Aliases
	}

	return config, nil
}

func (c *Config) merge(args Args) {
	if args.BroadcastAddress != "" {
		c.BroadcastAddress = args.BroadcastAddress
	}
	if args.Port != 0 {
		c.Port = args.Port
	}
}

// lookupAlias resolves an alias (case-insensitive) to its stored MAC address.
func (c Config) lookupAlias(alias string) (string, bool) {
	mac, ok := c.Aliases[strings.ToLower(alias)]
	return mac, ok
}

// saveConfig writes the given Config to disk at cfgPath, creating the directory if needed.
func saveConfig(config Config, cfgPath string) error {
	if err := os.MkdirAll(path.Dir(cfgPath), 0o755); err != nil {
		return fmt.Errorf("Error creating config directory: %w", err)
	}

	file, err := os.Create(cfgPath)
	if err != nil {
		return fmt.Errorf("Error opening config file for writing: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(config); err != nil {
		return fmt.Errorf("Error encoding config file: %w", err)
	}

	return nil
}
