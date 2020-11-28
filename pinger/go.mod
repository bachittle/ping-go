module github.com/bachittle/ping-go/pinger

go 1.15

require (
	github.com/bachittle/ping-go/utils v0.0.0-20201128190917-1230e44b481f
	github.com/jackpal/gateway v1.0.6
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace (
	github.com/bachittle/ping-go/utils => ../utils
	github.com/jackpal/gateway => ../ext/gateway
)
