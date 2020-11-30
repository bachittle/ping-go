package pinger

import (
	"github.com/bachittle/ping-go/utils"
	"testing"
)

// TestPinger tests pinger
func TestPinger(t *testing.T) {
	tests := []string{
		"localhost",
		"192.168.50.1", // change this number to your default gateway
		"google.com",
		"uwindsor.ca",
	}
	for _, test := range tests {
		t.Run(test, func(t *testing.T) {
			p := NewPinger()
			dst, err := utils.GetIPv4(test)
			if err != nil {
				t.Error("Error parsing IP:", err)
			}
			p.SetDst(dst)
			t.Log("testing", p.src, p.dst, p.amt)
			_, err = p.Ping()
			if err != nil {
				t.Error("Error in p.Ping():", err)
			}
		})
	}
}
