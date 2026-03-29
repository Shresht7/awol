package main

import (
	"fmt"
	"net"
	"os"
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
		os.Exit(1)
	}
	config.merge(args)

	switch args.SubCmd {
	case "list":
		if err := listAliases(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "alias":
		if err := setAlias(config, cfgPath, args.Rest); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "remove":
		if err := removeAlias(config, cfgPath, args.Rest); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		return
	case "config-path":
		fmt.Println(cfgPath)
		return
	case "help":
		showHelp()
		return
	case "version":
		showVersion()
		return
	default:
	}

	// No MAC address provided — show help
	if args.Mac == "" {
		showHelp()
		os.Exit(1)
	}

	// The argument is expected to be a MAC address or an alias defined in the config
	macAddress := args.Mac
	macAlias := ""

	// Check if the provided MAC address is an alias in the config
	if val, exists := config.lookupAlias(args.Mac); exists {
		macAddress = val
		macAlias = args.Mac
	}

	// Parse the MAC address
	mac, err := net.ParseMAC(macAddress)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing MAC Address [%s]: %v\n", macAddress, err)
		os.Exit(1)
	}

	// Create the magic packet
	magicPacket := makeMagicPacket(mac)

	// Send the magic packet via UDP broadcast (standard port for Wake-on-LAN is 9)
	sourceAddress := args.SourceAddress
	broadcastAddress := fmt.Sprintf("%s:%d", config.BroadcastAddress, config.Port)
	err = broadcastMagicPacket(sourceAddress, broadcastAddress, magicPacket)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	if macAlias != "" {
		fmt.Printf("Magic packet sent to %s [%s]\n", macAlias, macAddress)
	} else {
		fmt.Printf("Magic packet sent to %s\n", macAddress)
	}
}
