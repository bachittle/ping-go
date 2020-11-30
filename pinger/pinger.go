package pinger

import (
	"errors"
	"fmt"
	"github.com/bachittle/gateway"
	"github.com/bachittle/ping-go/utils"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"io"
	"net"
	"os"
	"time"
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

func (p Pinger) String() string {
	return fmt.Sprintf("Pinger{%v, %v, %v}", p.src, p.dst, p.amt)
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

// TimeoutError is for Ping when it times out.
// It returns a pinger object, and timeout interval (in milliseconds).
type TimeoutError struct {
	pinger  Pinger
	timeout int
}

func (e *TimeoutError) Error() string {
	return fmt.Sprint(e.pinger, " timed out after ", e.timeout, " milliseconds")
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

// Ping does the action of pinging a server. Returns the ICMP message response.
// to prevent a hanging ping, errors on a timeout number as a time (milliseconds)
//
// usage:
// - p.Ping() // standard ping with timeout, no output writing
// - p.Ping(100) // ping with timeout of 100 millisecond, no output writing
// - p.Ping(100, os.Stdout) // ping with timeout of 100 millisecond to stdout
func (p Pinger) Ping(args ...interface{}) ([]*icmp.Message, error) {
	timeout := 1000
	var ok bool
	var writer io.Writer
	for _, arg := range args {
		timeout, ok = arg.(int)
		if !ok {
			timeout = 1000
			writer, ok = arg.(io.Writer)
			if !ok {
				return nil, errors.New(fmt.Sprint("invalid arguments: ", args))
			}
		}
	}
	var msgList []*icmp.Message
	for i := 0; i < p.amt; i++ {
		conn, err := icmp.ListenPacket("ip:icmp", fmt.Sprint(p.src)) // packets from localhost
		if err != nil {
			return nil, err
		}
		defer conn.Close()
		chanMsg := make(chan *icmp.Message, 1)
		var msg *icmp.Message
		chanErr := make(chan error, 1)
		go func() {
			msg, err := p.PingOne(conn)
			if err != nil {
				chanErr <- err
			} else {
				chanMsg <- msg
			}
		}()
		select {
		case res := <-chanMsg:
			if writer != nil {
				fmt.Fprintln(writer, "got response ", res)
			}
			msg = res
		case res := <-chanErr:
			if writer != nil {
				fmt.Fprintln(writer, "got an error ", res)
			}
			err = res
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			err = &TimeoutError{p, timeout}
			if writer != nil {
				fmt.Fprintln(writer, "got a timeout error ", err)
			}
		}
		if err != nil {
			return msgList, err
		}
		msgList = append(msgList, msg)
	}
	return msgList, nil
}

// PingOne pings a server with one packet. Can also pass a connection as parameter.
func (p Pinger) PingOne(conn *icmp.PacketConn) (*icmp.Message, error) {
	var err error
	if conn == nil {
		conn, err = icmp.ListenPacket("ip:icmp", fmt.Sprint(p.src)) // packets from localhost
		if err != nil {
			return nil, err
		}
		defer conn.Close()
	}

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
		return nil, err
	}

	_, err = conn.WriteTo(reqBinary, &net.IPAddr{IP: p.dst, Zone: ""})
	if err != nil {
		return nil, err
	}

	respBinary := make([]byte, 1500)
	n, _, err := conn.ReadFrom(respBinary)
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
		//fmt.Printf("got %+v; want echo reply\n", respMsg)
	}
	return respMsg, nil
}
