package pinger

import (
	"fmt"
	"github.com/bachittle/ping-go/utils"
	"testing"
)

// TestPinger tests pinger
func TestPinger(t *testing.T) {
	tests := []string{
		"localhost",
		"192.168.50.10", // change this number to your default gateway
		"google.com",
		"uwindsor.ca", // known to block ICMP packets
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
			err = p.Ping()
			if err != nil {
				_, ok := err.(*TimeoutError)
				if !ok {
					t.Error("Error in p.Ping():", err)
				} else {
					// timeout error isn't necessarily a coding bug, log just in case
					t.Log(fmt.Sprint(
						"\n---------------TIMEOUT---------------\n",
						err,
						".\nServer is either blocking ICMP packets or your internet is down. ",
						"\n-------------------------------------\n",
					))
				}
			}
		})
	}
}
