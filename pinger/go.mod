module github.com/bachittle/ping-go/pinger

go 1.15

require (
	github.com/bachittle/gateway v1.0.8
	github.com/bachittle/ping-go/utils v0.0.0-20201128190917-1230e44b481f
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace (
	github.com/bachittle/ping-go/utils => ../utils
)
