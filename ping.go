package main

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"net"
	"os"
	"runtime"
)

func main() {
	// analyze runtime OS
	fmt.Println(runtime.GOOS)

	conn, err := icmp.ListenPacket("ip:icmp", "127.0.0.1")
	if err != nil {
		fmt.Println("ERROR:", err.Error())
	}
	defer conn.Close()

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
		fmt.Println(err.Error())
	}

	_, err = conn.WriteTo(reqBinary, &net.IPAddr{IP: net.IPv4(127, 0, 0, 1), Zone: "en0"})
	if err != nil {
		fmt.Println(err.Error())
	}

	respBinary := make([]byte, 1500)
	n, peer, err := conn.ReadFrom(respBinary)
	if err != nil {
		fmt.Println(err)
	}
	respMessage, err := icmp.ParseMessage(1, respBinary[:n])
	if err != nil {
		fmt.Println(err)
	}
	switch respMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("got reflection from %v", peer)
	default:
		fmt.Printf("got %+v; want echo reply", respMessage)
		fmt.Println(respBinary[:n])
	}
}
