package utils

import (
	"errors"
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

// GetIPv4 takes a generic string and gets a single IPv4 address of type net.IP.
// if no IPv4 addresses exist return an error.
func GetIPv4(hostname string) (net.IP, error) {
	ipArr, err := net.LookupIP(hostname)
	if err != nil {
		return nil, err
	}
	for _, ip := range ipArr {
		if ip.To4() != nil {
			return ip, nil
		}
	}
	return nil, errors.New("No IPv4 addresses found")
}

// GetIPv4CIDR returns an array of IPs from a CIDR hostname
// EX: "192.168.2.1/24" will get 256 domains from 192.168.2.0 - 192.168.2.255
func GetIPv4CIDR(hostnameCIDR string) (IPs []net.IP, err error) {
	// check if hostname string is valid first
	i := strings.Index(hostnameCIDR, "/")
	if i == -1 {
		err = errors.New("Invalid index")
		return
	}
	ipCIDR, ipNet, err := net.ParseCIDR(hostnameCIDR)
	if err != nil {
		// try to look up IP address first
		var newIPs []string
		newIPs, err = net.LookupHost(hostnameCIDR[:i])
		for _, ip := range newIPs {
			ipCIDR, ipNet, err = net.ParseCIDR(fmt.Sprint(ip, hostnameCIDR[i:]))
			if err == nil {
				break
			}
		}
	}
	// fmt.Println("ip:", ipCIDR)
	// fmt.Println("ipNet:", ipNet)
	for ip := ipCIDR.Mask(ipNet.Mask); ipNet.Contains(ip); IncrementIP(ip) {
		newIP := make(net.IP, len(ip))
		copy(newIP, ip)
		IPs = append(IPs, newIP)
	}
	// error checking the lengths before return
	var CIDRnum int

	CIDRnum, err = strconv.Atoi(hostnameCIDR[i+1:])
	if err != nil {
		err = errors.New(fmt.Sprint("Invalid number:", hostnameCIDR[i+1:]))
		return
	}
	CIDRlen := int(math.Pow(2, float64(32-CIDRnum)))
	if CIDRlen != len(IPs) {
		err = errors.New(fmt.Sprint("Error! Expected length", CIDRlen, ", got length", len(IPs)))
		return
	}
	return
}

// IncrementIP will increment a given IP by doing byte manipulation, then return the result.
func IncrementIP(ip net.IP) {
	// fmt.Printf("ip: %#v\n", ip)
	var i int
	for i = len(ip) - 1; i >= 0; i-- {
		// fmt.Printf("increment this ip: %#v\n", ip[i]+1)
		if ip[i] != 0xff {
			ip[i]++
			break
		}
	}
	if i == 0 && ip[i] == 0xff {
		ip = nil
	}
}

/*

// don't need this as now I am using github.com/jackpal/gateway.
// Will save this for a rainy day


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
