# `awol`

A simple command-line-interface to dispatch Wake-on-LAN magic packets!

This allows you to ***remotely*** switch on a computer/machine on your local network, (as long as it supports Wake-on-LAN and is configured to do so).

## Usage

Use `--help` or `-h` to see the help message.

To wake a computer, provide its MAC address as an argument to the `awol` command. For example:

```bash
awol <mac_address>
```
Example:

```bash
awol A1:2B:C3:4D:5E:F7
```

## Installation

```bash
go install github.com/Shresht7/awol@latest
```

## Wake-on-LAN

Wake-on-LAN (WoL) is a networking standard that allows a computer to be turned on or awakened from a low power state remotely. This is achieved by sending a specially crafted network packet, known as a _"magic packet"_ to the target computer's network interface. The magic packet contains the target computer's MAC address repeated multiple (16) times, which allows the network interface to recognize it and trigger the power-on process.

See [Wake-on-LAN - Wikipedia](https://en.wikipedia.org/wiki/Wake-on-LAN) for more details.

### Configuration

> [!IMPORTANT]
> To use Wake-on-LAN, the target computer must support it and be configured to allow it. This typically involves enabling WoL in the computer's **BIOS/UEFI** settings and ensuring that the network interface is set to allow wake-up events. Additionally, the computer must be connected to a power source and the network interface must be active (not in a completely powered-off state) for WoL to work.

---

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
