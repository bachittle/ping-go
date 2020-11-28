package pinger

import (
	"github.com/bachittle/ping-go/utils"
	"testing"
)

// TestPinger
func TestPinger(t *testing.T) {
	p := NewPinger()
	dst, err := utils.GetIPv4("192.168.50.1")
	if err != nil {
		t.Error("Error parsing IP:", err)
	}
	p.dst = dst
	err = p.Ping()
	if err != nil {
		t.Error("Error in p.Ping():", err)
	}
}
