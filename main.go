package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

// ----
// MAIN
// ----

// The main entrypoint of the application
func main() {

	// Parse the command-line arguments
	args := parseCommandLineArgs()

	// Check for help flag and show the help message
	if args.Help {
		showHelp()
		return
	}

	// Check for version flag and show the version information
	if args.Version {
		showVersion()
		return
	}

	// Load the configuration file
	cfgPath := getConfigPath()
	config, err := readConfig(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config: %v\n", err)
		return
	}
	config.merge(args)

	switch args.SubCmd {
	case "list":
		listAliases(config)
		return
	case "alias":
		setAlias(config, cfgPath, args.Rest)
		return
	case "remove":
		removeAlias(config, cfgPath, args.Rest)
		return
	case "help":
		showHelp()
		return
	case "version":
		showVersion()
		return
	default:
	}

	// The argument is expected to be a MAC address or an alias defined in the config
	macAddress := args.Mac
	macAlias := ""

	// Check if the provided MAC address is an alias in the config
	if val, exists := config.Aliases[strings.ToLower(args.Mac)]; exists {
		macAddress = val
		macAlias = args.Mac
	}

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Address [%s]: %v\n", macAddress, err)
		return
	}

	// Create the magic packet
	magicPacket := makeMagicPacket(mac)

	// Send the magic packet via UDP broadcast (standard port for Wake-on-LAN is 9)
	broadcastAddress := fmt.Sprintf("%s:%d", "255.255.255.255", config.Port)
	err = broadcastMagicPacket(broadcastAddress, magicPacket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		return
	}

	if macAlias != "" {
		fmt.Printf("Magic packet sent to %s [%s]\n", macAlias, macAddress)
	} else {
		fmt.Printf("Magic packet sent to %s\n", macAddress)
	}
}
