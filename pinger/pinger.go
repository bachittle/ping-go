package pinger

import (
	"errors"
	"fmt"
	"github.com/bachittle/gateway"
	"github.com/bachittle/ping-go/utils"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"math/rand"
	"net"
	"time"
)

// Pinger structure has the following properties
// - src (source) IP address (most likely 127.0.0.1)
// - dst (destination) IP address (set by the user)
// - amt (amount) of pings to the dst IP address.
type Pinger struct {
	src  net.IP
	dst  net.IP
	conn *icmp.PacketConn
	amt  int
}

func (p Pinger) String() string {
	return fmt.Sprintf("Pinger{%v, %v, %v}", p.src, p.dst, p.amt)
}

// NewPinger creates a default Pinger by calling p.Default() and returns p
func NewPinger() Pinger {
	p := Pinger{nil, nil, nil, 0}
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

	if src != nil {
		p.src = src
	} else {
		p.src = defaultIP
	}

	p.SetDst(dst)

	if amt != nil {
		p.amt = *amt
	} else {
		p.amt = defaultAmt
	}
	return nil
}

// SetAmt to set private value amt
func (p *Pinger) SetAmt(amt int) int {
	p.amt = amt
	return amt
}

// SetSrc to set private value src
func (p *Pinger) SetSrc(src net.IP) (net.IP, error) {
	if src == nil {
		return nil, errors.New("src must not be nil")
	}
	p.src = src
	return p.src, nil
}

// NewConn creates the icmp packet "connection"
func (p Pinger) NewConn() (*icmp.PacketConn, error) {
	return icmp.ListenPacket("ip:icmp", fmt.Sprint(p.src)) // packets from localhost
}

// TimeoutError is for Ping when it times out.
// It returns a pinger object, and timeout interval (in milliseconds).
type TimeoutError struct {
	Pinger  Pinger
	Timeout int
}

func (e *TimeoutError) Error() string {
	return fmt.Sprint(e.Pinger, " timed out after ", e.Timeout, " milliseconds")
}

// SetDst is a setter function that does some required changes while setting dst,
// 		including changing the src IP
func (p *Pinger) SetDst(dst net.IP) (net.IP, error) {
	if dst == nil {
		return nil, errors.New("dst must not be nil")
	}
	p.dst = dst
	localhostIP, err := utils.GetIPv4("localhost")
	if err != nil {
		return nil, err
	}
	/*
		fmt.Println("src:", p.src)
		fmt.Println("dst:", p.dst)
		fmt.Println("localhost:", localhostIP)
		fmt.Println("equality:", p.dst.Equal(localhostIP))
	*/
	if p.dst.Equal(localhostIP) {
		p.src = p.dst
	}
	return p.dst, nil
}

// PingPong sends an echo request with Ping and receives a result with Pong.
// usage:
// p.PingPong() 	// uses default timeout time of 1000ms
// p.PingPong(500) // set the timeout to 500ms
func (p Pinger) PingPong(args ...interface{}) ([]*icmp.Message, error) {
	timeout := 1000
	for _, arg := range args {
		temp, ok := arg.(int)
		if ok {
			timeout = temp
			continue
		}
	}
	err := p.Ping()
	if err != nil {
		return nil, err
	}
	return p.Pong(timeout)
}

// Ping sends a ping to the specified dst with amt packets
func (p *Pinger) Ping() error {
	conn, err := p.NewConn()
	if err != nil {
		return err
	}
	p.conn = conn
	for i := 0; i < p.amt; i++ {
		err = p.SendOnePing(i, conn)
	}
	return err
}

// Pong receives a ping from the specified dst with amt packets asynchronously with timeout
func (p Pinger) Pong(timeout int) (msgList []*icmp.Message, err error) {
	cErr := make(chan error, 1)
	cMsg := make(chan *icmp.Message, 1)
	for i := 0; i < p.amt; i++ {
		go func() {
			msg, err := p.RecvOnePong()
			if err != nil {
				cErr <- err
				return
			}
			cMsg <- msg
		}()
	}
	for i := 0; i < p.amt; i++ {
		select {
		case res := <-cErr:
			err = res
		case res := <-cMsg:
			msgList = append(msgList, res)
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			err = errors.New("timeout")
		}
	}
	return
}

// SendOnePing pings a server with one packet. Can also pass a connection as parameter.
func (p Pinger) SendOnePing(seq int, conn *icmp.PacketConn) error {
	var err error
	if conn == nil {
		conn, err = p.NewConn()
		if err != nil {
			return err
		}
		p.conn = conn
	}

	reqMsg := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID:   rand.Intn(65535),
			Seq:  seq,
			Data: []byte(""),
		},
	}

	reqBinary, err := reqMsg.Marshal(nil)
	if err != nil {
		return err
	}

	ipAddr := &net.IPAddr{IP: p.dst, Zone: ""}
	_, err = conn.WriteTo(reqBinary, ipAddr)
	return err
}

// RecvOnePong receives the result of a SendPing message. Must include packet connection.
func (p Pinger) RecvOnePong() (*icmp.Message, error) {
	if p.conn == nil {
		return nil, errors.New("no conn, cannot read pong")
	}
	respBinary := make([]byte, 1500)
	n, _, err := p.conn.ReadFrom(respBinary)
	if err != nil {
		return nil, err
	}
	respMsg, err := icmp.ParseMessage(1, respBinary[:n])
	if err != nil {
		return nil, err
	}
	switch respMsg.Type {
	case ipv4.ICMPTypeEchoReply:
		//fmt.Printf("got reflection from %v\n", peer)
	default:
		// fmt.Printf("got %+v; want echo reply\n", respMsg)
		return respMsg, errors.New("notEchoReply")
	}
	return respMsg, nil
}
