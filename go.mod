module github.com/bachittle/ping-go

go 1.15

require (
	github.com/bachittle/ping-go/pinger v0.0.0-20201204220814-575c884bcb87
	github.com/bachittle/ping-go/utils v0.0.0-20201204220814-575c884bcb87
)

replace (
	github.com/bachittle/ping-go/pinger => ./pinger
	github.com/bachittle/ping-go/utils => ./utils
)
