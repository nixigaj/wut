package main

import (
	"fmt"
	"github.com/nixigaj/go-default-route"
	"net"
)

func getDefaultIPv6Address() (string, error) {
	interf, err := defaultroute.DefaultRouteInterface()
	if err != nil {
		return "", err
	}
	// Get the default network interface

	// Get the addresses associated with the default interface
	addrs, err := interf.Addrs()
	if err != nil {
		return "", err
	}

	// Iterate through the addresses and find the first IPv6 address
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok {
			if ipNet.IP.To16() != nil && ipNet.IP.To4() == nil {
				return ipNet.IP.String(), nil
			}
		}
	}

	return "", fmt.Errorf("No IPv6 address found for the default interface")
}

func main() {
	ipv6Address, err := getDefaultIPv6Address()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Printf("Default IPv6 Address: %s\n", ipv6Address)
}
