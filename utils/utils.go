package utils

import (
	"errors"
	"fmt"
	"net"
)

// GetIPv4 takes a generic string and gets a single IPv4 address of type net.IP.
// if no IPv4 addresses exist return an error.
func GetIPv4(hostname string) (net.IP, error) {
	ipArr, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Println(err)
	}
	for _, ip := range ipArr {
		if ip.To4() != nil {
			return ip, nil
		}
	}
	return nil, errors.New("No IPv4 addresses found")
}
