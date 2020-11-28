package utils

import (
	"errors"
	"fmt"
	"net"
	//"github.com/jackpal/gateway"
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

/*

don't need this as now I am using github.com/jackpal/gateway.
Will save this for a rainy day


// LocalIP has an IP and an interface associated with the IP.
type LocalIP struct {
	ip    net.IP
	iface string
}

// GetLocalIPs returns your IP for all interfaces on the machine as an array of the form {ip, interface} (type LocalIP)
func GetLocalIPs() ([]net.IP, error) {
	var ips []net.IP
	ifaces, err := net.Interfaces()
	// handle err
	if err != nil {
		return ips, err
	}
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		// handle err
		if err != nil {
			return ips, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			ips = append(ips, ip)
		}
	}
	return ips, nil
}

// GetLocalIPv4 only gets IPv4 addresses, used in pinger at the moment (no ipv6)
func GetLocalIPv4() ([]net.IP, error) {
	var newIPs []net.IP
	ips, err := GetLocalIPs()
	if err != nil {
		return ips, err
	}
	for _, ip := range ips {
		if ip.To4() == nil {
			continue
		}
		newIPs = append(newIPs, ip)
	}
	return newIPs, nil
}

// GetGatewayIPv4 should get the IP address needed for ICMP packets
func GetGatewayIPv4() (net.IP, error) {
	IPs, err := GetLocalIPv4()
	if err != nil {
		return nil, err
	}
	gateIP, err := gateway.DiscoverInterface()
	if err != nil {
		return nil, err
	}
	fmt.Println(gateIP, IPs)
	return nil, nil
}
*/
