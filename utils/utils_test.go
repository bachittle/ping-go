package utils

import (
	"testing"
)

// TestGetIPv4 tests an array of hostname strings to see if GetIPv4 returns a net.IP that is of type IPv4.
func TestGetIPv4(t *testing.T) {
	passTests := []string{
		"127.0.0.1",
		"localhost",
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
