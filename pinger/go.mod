module github.com/bachittle/ping-go/pinger


go 1.15

require (
	github.com/bachittle/ping-go/utils v0.0.0-00010101000000-000000000000
	github.com/jackpal/gateway v1.0.6
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace (
	github.com/jackpal/gateway => ../ext/gateway
	github.com/bachittle/ping-go/utils => ../utils
)
