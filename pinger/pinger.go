package pinger

import (
	"errors"
	"fmt"
	"github.com/bachittle/ping-go/utils"
	"github.com/jackpal/gateway"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
)

// Pinger structure has the following properties
// - src (source) IP address (most likely 127.0.0.1)
// - dst (destination) IP address (set by the user)
// - amt (amount) of pings to the dst IP address.
type Pinger struct {
	src net.IP
	dst net.IP
	amt int
}

// NewPinger creates a default Pinger by calling p.Default() and returns p
func NewPinger() Pinger {
	p := Pinger{nil, nil, 0}
	p.Default(nil, nil, nil)
	return p
}

// Default values for ping if the user is lazy and does not want to specify details
// if nil, do default value.
//
// Default values for each field:
// - src: Default Gateway Interface (192.168.2.*)
// - dst: Default Gateway Interface (192.168.2.*)
// - amt: 32
func (p *Pinger) Default(src net.IP, dst net.IP, amt *int) error {

	defaultIP, err := gateway.DiscoverInterface()
	if err != nil {
		return err
	}
	defaultAmt := 32

	if dst != nil {
		p.dst = dst
	} else {
		p.dst = defaultIP
	}
	if src != nil {
		p.src = src
	} else {
		p.src = defaultIP
	}

	if amt != nil {
		p.amt = *amt
	} else {
		p.amt = defaultAmt
	}
	return nil
}

// SetDst is a setter function that does some required changes while setting dst,
// 		including changing the src IP
func (p *Pinger) SetDst(dst net.IP) error {
	if dst == nil {
		return errors.New("destination IP must not be nil")
	}
	p.dst = dst
	localhostIP, err := utils.GetIPv4("localhost")
	if err != nil {
		return err
	}
	fmt.Println("src:", p.src)
	fmt.Println("dst:", p.dst)
	fmt.Println("localhost:", localhostIP)
	fmt.Println("equality:", p.dst.Equal(localhostIP))
	if p.dst.Equal(localhostIP) {
		p.src = p.dst
	}
	return nil
}

// Ping does the action of pinging a server.
func (p Pinger) Ping() error {
	conn, err := icmp.ListenPacket("ip:icmp", fmt.Sprint(p.src)) // packets from localhost
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < p.amt; i++ {
		reqMsg := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID:   os.Getpid() & 0xffff,
				Seq:  1,
				Data: []byte(""),
			},
		}

		reqBinary, err := reqMsg.Marshal(nil)
		if err != nil {
			return err
		}

		_, err = conn.WriteTo(reqBinary, &net.IPAddr{IP: p.dst, Zone: ""})
		if err != nil {
			return err
		}

		respBinary := make([]byte, 1500)
		n, peer, err := conn.ReadFrom(respBinary)
		if err != nil {
			return err
		}
		respMessage, err := icmp.ParseMessage(1, respBinary[:n])
		if err != nil {
			return err
		}
		switch respMessage.Type {
		case ipv4.ICMPTypeEchoReply:
			fmt.Printf("got reflection from %v\n", peer)
		default:
			fmt.Printf("got %+v; want echo reply\n", respMessage)
		}
	}
	return nil
}
