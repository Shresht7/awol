package main

import (
	"fmt"
	"net"
	"os"
	"sort"
	"strings"
)

// --------
// COMMANDS
// --------

func listAliases(config Config) {
	if len(config.Aliases) == 0 {
		fmt.Println("No aliases defined in the config file.")
		return
	}

	keys := make([]string, 0, len(config.Aliases))
	for alias := range config.Aliases {
		keys = append(keys, alias)
	}
	sort.Strings(keys)

	for _, alias := range keys {
		fmt.Printf("%s\t%s\n", alias, config.Aliases[alias])
	}
}

func setAlias(config Config, cfgPath string, rest []string) {
	if len(rest) < 2 {
		fmt.Fprintln(os.Stderr, "Error: Alias and MAC address must be provided for the 'alias' command.")
		return
	}
	alias, mac := rest[0], rest[1]

	// Validate the MAC address
	parsed, err := net.ParseMAC(mac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Address [%s]: %v\n", mac, err)
		return
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
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		return
	}

	fmt.Printf("Alias '%s' set to MAC address '%s'\n", normalizedAlias, normalizedMAC)
}

func removeAlias(config Config, cfgPath string, rest []string) {
	if len(rest) < 1 {
		fmt.Fprintln(os.Stderr, "Error: Alias must be provided for the 'remove' command.")
		return
	}
	alias := rest[0]

	// Check if the alias exists in the config
	if _, exists := config.Aliases[strings.ToLower(alias)]; !exists {
		fmt.Fprintf(os.Stderr, "Error: Alias '%s' does not exist in the config file.\n", alias)
		return
	}

	// Remove the alias from the config
	delete(config.Aliases, strings.ToLower(alias))

	if err := saveConfig(config, cfgPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		return
	}

	fmt.Printf("Alias '%s' removed from the config file.\n", alias)
}
