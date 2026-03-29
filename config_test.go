package main

import (
	"os"
	"reflect"
	"testing"
)

func TestGetConfigPath(t *testing.T) {
	// Test default config path
	defaultPath := getConfigPath()

	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home directory: %v", err)
	}

	expectedDefaultPath := home + "/.config/awol/config.json"
	if defaultPath != expectedDefaultPath {
		t.Errorf("Expected default config path '%s', got '%s'", expectedDefaultPath, defaultPath)
	}

	// Test custom config path via environment variable
	customPath := "/custom/config/path.json"
	os.Setenv("AWOL_CONFIG_PATH", customPath)
	defer os.Unsetenv("AWOL_CONFIG_PATH")

	if getConfigPath() != customPath {
		t.Errorf("Expected custom config path '%s', got '%s'", customPath, getConfigPath())
	}

	// Test empty environment variable
	os.Setenv("AWOL_CONFIG_PATH", "")
	if getConfigPath() != expectedDefaultPath {
		t.Errorf("Expected default config path '%s' when environment variable is empty, got '%s'", expectedDefaultPath, getConfigPath())
	}
}

func TestReadConfig(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "awol_config_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write a sample config to the temp file
	sampleConfig := `{
		"broadcast": "255.255.255.255",
		"port": 7,
		"aliases": {
			"laptop": "00:11:22:33:44:55",
			"desktop": "66:77:88:99:AA:BB"
		}
	}`
	if _, err := tmpFile.Write([]byte(sampleConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read the config using the readConfig function
	config, err := readConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Validate the config values
	if config.BroadcastAddress != "255.255.255.255" {
		t.Errorf("Expected broadcast address '255.255.255.255', got '%s'", config.BroadcastAddress)
	}
	if config.Port != 7 {
		t.Errorf("Expected port 7, got %d", config.Port)
	}
	if config.Aliases["laptop"] != "00:11:22:33:44:55" {
		t.Errorf("Expected alias 'laptop' to be '00:11:22:33:44:55', got '%s'", config.Aliases["laptop"])
	}
	if config.Aliases["desktop"] != "66:77:88:99:AA:BB" {
		t.Errorf("Expected alias 'desktop' to be '66:77:88:99:AA:BB', got '%s'", config.Aliases["desktop"])
	}

	// Clean up the temporary file
	if err := os.Remove(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}

	// Test reading a non-existent config file
	_, err = readConfig("non_existent_config.json")
	if err != nil {
		t.Fatalf("Expected no error when reading non-existent config file, got: %v", err)
	}
}

func TestReadConfigWithMissingFields(t *testing.T) {
	// Create a temporary config file with missing fields
	tmpFile, err := os.CreateTemp("", "awol_config_missing_fields_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Write a sample config with missing fields to the temp file
	sampleConfig := `{
		"port": 7
	}`
	if _, err := tmpFile.Write([]byte(sampleConfig)); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read the config using the readConfig function
	config, err := readConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Validate that missing fields are set to default values
	if config.BroadcastAddress != "255.255.255.255" {
		t.Errorf("Expected default broadcast address '255.255.255.255', got '%s'", config.BroadcastAddress)
	}
	if config.Port != 7 {
		t.Errorf("Expected port 7, got %d", config.Port)
	}

	// Clean up the temporary file
	if err := os.Remove(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

func TestSaveReadRoundTrip(t *testing.T) {
	// Create a temporary config file
	tmpFile, err := os.CreateTemp("", "awol_config_round_trip_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	// Create a sample config
	originalConfig := Config{
		BroadcastAddress: "255.255.255.255",
		Port:             7,
		Aliases: map[string]string{
			"laptop":  "00:11:22:33:44:55",
			"desktop": "66:77:88:99:AA:BB",
		},
	}

	// Save the config to the temp file
	if err := saveConfig(originalConfig, tmpFile.Name()); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temp file: %v", err)
	}

	// Read the config back from the temp file
	readConfig, err := readConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}

	// Compare the original and read configs
	if !reflect.DeepEqual(originalConfig, readConfig) {
		t.Errorf("Expected config %+v, got %+v", originalConfig, readConfig)
	}

	// Clean up the temporary file
	if err := os.Remove(tmpFile.Name()); err != nil {
		t.Fatalf("Failed to remove temp file: %v", err)
	}
}

func TestLookupAlias(t *testing.T) {
	// Create a sample config with aliases
	config := Config{
		BroadcastAddress: "255.255.255.255",
		Port:             7,
		Aliases: map[string]string{
			"laptop":  "00:11:22:33:44:55",
			"desktop": "66:77:88:99:AA:BB",
		},
	}

	// Test alias lookup (case-insensitive)
	if mac, ok := config.lookupAlias("laptop"); !ok || mac != "00:11:22:33:44:55" {
		t.Errorf("Expected alias 'laptop' to resolve to '00:11:22:33:44:55', got '%s'", mac)
	}
	if mac, ok := config.lookupAlias("DESKTOP"); !ok || mac != "66:77:88:99:AA:BB" {
		t.Errorf("Expected alias 'DESKTOP' to resolve to '66:77:88:99:AA:BB', got '%s'", mac)
	}
	if mac, ok := config.lookupAlias("desktop"); !ok || mac != "66:77:88:99:AA:BB" {
		t.Errorf("Expected alias 'desktop' to resolve to '66:77:88:99:AA:BB', got '%s'", mac)
	}
	if mac, ok := config.lookupAlias("unknown"); ok || mac != "" {
		t.Errorf("Expected alias 'unknown' to resolve to '', got '%s'", mac)
	}
}

func TestMerge(t *testing.T) {
	// Create a sample config
	config := Config{
		BroadcastAddress: "255.255.255.255",
		Port:             7,
		Aliases: map[string]string{
			"laptop":  "00:11:22:33:44:55",
			"desktop": "66:77:88:99:AA:BB",
		},
	}
	// Create sample args with overrides
	args := Args{
		BroadcastAddress: "192.168.1.255",
		Port:             9,
	}

	// Merge the args into the config
	config.merge(args)

	// Validate that the config fields were updated
	if config.BroadcastAddress != "192.168.1.255" {
		t.Errorf("Expected broadcast address '192.168.1.255', got '%s'", config.BroadcastAddress)
	}
	if config.Port != 9 {
		t.Errorf("Expected port 9, got %d", config.Port)
	}
	// Validate that aliases were not modified
	expectedAliases := map[string]string{
		"laptop":  "00:11:22:33:44:55",
		"desktop": "66:77:88:99:AA:BB",
	}
	if !reflect.DeepEqual(config.Aliases, expectedAliases) {
		t.Errorf("Expected aliases %+v, got %+v", expectedAliases, config.Aliases)
	}
}
