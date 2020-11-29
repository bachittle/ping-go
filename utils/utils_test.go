package utils

import (
	//"github.com/bachittle/gateway"
	//"net"
	"fmt"
	"testing"
)

// TestGetIPv4 tests an array of hostname strings to see if GetIPv4 returns a net.IP that is of type IPv4.
func TestGetIPv4(t *testing.T) {
	passTests := []string{
		"127.0.0.1",
		"localhost",
		"192.168.50.1", // set this to your default gateway
		"google.com",
	}

	failTests := []string{
		"::1",
	}
	// These tests should NOT get an error.
	// If they do, they fail the test
	t.Log("testing cases that pass")
	for _, test := range passTests {
		t.Run(test, func(t *testing.T) {
			ip, err := GetIPv4(test)
			if err != nil {
				t.Error("Got an error:", err)
			} else {
				t.Log("Got a response:", ip)
			}
		})
	}

	// These tests SHOULD get an error
	// if they don't, then fail the test
	t.Log("testing cases that fail")
	for _, test := range failTests {
		t.Run(test, func(t *testing.T) {
			ip, err := GetIPv4(test)
			if err != nil {
				t.Log("Got an error:", err)
			} else {
				t.Error("Got a response:", ip)
			}
		})
	}
}

func TestGetIPv4CIDR(t *testing.T) {
	passTests := []string{
		"192.168.50.1/24", // set this to your default gateway
		"google.com/24",
	}
	for _, test := range passTests {
		t.Run(test, func(t *testing.T) {
			IPs, err := GetIPv4CIDR(test)
			if err != nil {
				t.Error(err)
			}
			t.Log(fmt.Sprint("found ", len(IPs), " IP addresses for ", test))
		})
	}
}

/*
func TestGetLocalIPs(t *testing.T) {
	// IP addresses set by the computer. Set this manually
	// Go to command line
	//  	- windows: type 'arp -a'
	// 		- linux:   type 'ifconfig' or 'ip addr'
	myLocalIPs := []string{
		"192.168.50.76",
	}
	//t.Log(GetGatewayIPv4())
	for _, test := range myLocalIPs {
		t.Run(test, func(t *testing.T) {
			cmp, err := gateway.DiscoverInterface()
			if err != nil {
				t.Error(err)
			}
			if !cmp.Equal(net.ParseIP(test)) {
				t.Error(cmp, "!=", test)
			}

			ips, err := GetLocalIPv4()
			if err != nil {
				t.Error(err)
			}
			t.Log(ips)
			for _, ip := range ips {
				if fmt.Sprint(ip) == test {
					return
				}
			}
			t.Error("Could not find local IP")
		})
	}
}
*/
