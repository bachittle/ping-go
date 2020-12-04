package main

import (
	"fmt"
	"github.com/bachittle/ping-go/pinger"
	"github.com/bachittle/ping-go/utils"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]
	var err error
	var ips []net.IP

	// must have at least a destination IP address. If not, print usage.
	if len(args) == 0 {
		fmt.Println("usage: ping-go [-c count] destination")
		os.Exit(0)
	}
	p := pinger.NewPinger()
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-c":
			i++ // use next argument i=i+1
			amt, err := strconv.Atoi(args[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			p.SetAmt(amt)
			break
		default:
			// use this argument i=i
			if strings.Contains(args[i], "/") {
				ips, err = utils.GetIPv4CIDR(args[i])
			} else {
				var ip net.IP
				ip, err = utils.GetIPv4(args[i])
				ips = append(ips, ip)
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
		for _, ip := range ips {
			p.SetDst(ip)
			msgs, err := p.PingPong()
			if err != nil {
				fmt.Println("error:", err)
			}
			if len(msgs) > 0 {
				fmt.Println("got result!", msgs[0])
			}
		}
	}
}
