package pinger

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"reflect"
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
// - src: 127.0.0.1
// - dst: 127.0.0.1
// - amt: 32
func (p *Pinger) Default(src net.IP, dst net.IP, amt *int) {
	defaultIP := net.ParseIP("127.0.0.1")
	defaultAmt := 32

	if src != nil {
		p.src = src
	} else {
		p.src = defaultIP
	}
	if dst != nil {
		p.dst = dst
	} else {
		p.dst = defaultIP
	}
	if amt != nil {
		p.amt = *amt
	} else {
		p.amt = defaultAmt
	}
}

// Ping does the action of pinging a server.
func (p Pinger) Ping() error {
	conn, err := icmp.ListenPacket("ip:icmp", fmt.Sprint(p.src)) // packets from localhost
	if err != nil {
		fmt.Println("ERROR:", err.Error())
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

		fmt.Println(reflect.TypeOf(conn))
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
