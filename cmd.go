package main

import (
	"fmt"
	"net"
	"sort"
	"strings"
)

// --------
// COMMANDS
// --------

func listAliases(config Config) error {
	if len(config.Aliases) == 0 {
		fmt.Println("No aliases defined in the config file.")
		return nil
	}

	keys := make([]string, 0, len(config.Aliases))
	for alias := range config.Aliases {
		keys = append(keys, alias)
	}
	sort.Strings(keys)

	for _, alias := range keys {
		fmt.Printf("%s\t%s\n", alias, config.Aliases[alias])
	}
	return nil
}

func setAlias(config Config, cfgPath string, rest []string) error {
	if len(rest) < 2 {
		return fmt.Errorf("Alias and MAC address must be provided for the 'alias' command")
	}
	alias, mac := rest[0], rest[1]

	// Validate the MAC address
	parsed, err := net.ParseMAC(mac)
	if err != nil {
		return fmt.Errorf("Error parsing MAC Address [%s]: %w", mac, err)
	}

	// Normalize: aliases are lowercase, MACs are uppercase
	normalizedAlias := strings.ToLower(alias)
	normalizedMAC := strings.ToUpper(parsed.String())

	// Update the config with the new alias
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}

	config.Aliases[normalizedAlias] = normalizedMAC

	if err := saveConfig(config, cfgPath); err != nil {
		return err
	}

	fmt.Printf("Alias '%s' set to MAC address '%s'\n", normalizedAlias, normalizedMAC)
	return nil
}

func removeAlias(config Config, cfgPath string, rest []string) error {
	if len(rest) < 1 {
		return fmt.Errorf("Alias must be provided for the 'remove' command")
	}
	alias := rest[0]

	// Check if the alias exists in the config
	if _, exists := config.lookupAlias(alias); !exists {
		return fmt.Errorf("Alias '%s' does not exist in the config file", alias)
	}

	// Remove the alias from the config
	delete(config.Aliases, strings.ToLower(alias))

	if err := saveConfig(config, cfgPath); err != nil {
		return err
	}

	fmt.Printf("Alias '%s' removed from the config file.\n", alias)
	return nil
}
