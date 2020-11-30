package main

import (
	"fmt"
	"github.com/bachittle/ping-go/pinger"
	"github.com/bachittle/ping-go/utils"
	"os"
	"strconv"
)

func main() {
	args := os.Args[1:]
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
			ip, err := utils.GetIPv4(args[i])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			p.SetDst(ip)
			break
		}
	}
	_, err := p.Ping(os.Stdout)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
