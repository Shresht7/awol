package main

import (
	"flag"
	"fmt"
	"net"
	"os"
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

	for alias, mac := range config.Aliases {
		fmt.Printf("%s\t%s\n", alias, mac)
	}
}

func setAlias(config Config, cfgPath string) {
	alias := flag.Arg(1)
	mac := flag.Arg(2)

	if alias == "" || mac == "" {
		fmt.Fprintln(os.Stderr, "Error: Alias and MAC address must be provided for the 'alias' command.")
		return
	}

	// Validate the MAC address
	_, err := net.ParseMAC(mac)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Address [%s]: %v\n", mac, err)
		return
	}

	// Update the config with the new alias
	if config.Aliases == nil {
		config.Aliases = make(map[string]string)
	}

	config.Aliases[strings.ToLower(alias)] = mac

	if err := saveConfig(config, cfgPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		return
	}

	fmt.Printf("Alias '%s' set to MAC address '%s'\n", alias, mac)
}

func removeAlias(config Config, cfgPath string) {
	alias := flag.Arg(1)

	if alias == "" {
		fmt.Fprintln(os.Stderr, "Error: Alias must be provided for the 'remove' command.")
		return
	}

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
