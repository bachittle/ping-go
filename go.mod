module github.com/bachittle/ping-go

go 1.15

require (
	github.com/bachittle/ping-go/pinger v0.0.0-20201128190917-1230e44b481f
	github.com/bachittle/ping-go/utils v0.0.0-20201128190917-1230e44b481f
	github.com/jackpal/gateway v1.0.6 // indirect
	golang.org/x/net v0.0.0-20201110031124-69a78807bb2b
)

replace (
	github.com/bachittle/ping-go/pinger => ./pinger
	github.com/bachittle/ping-go/utils => ./utils
	github.com/jackpal/gateway => ./ext/gateway
)
