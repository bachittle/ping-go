module github.com/bachittle/ping-go

go 1.15

require (
	github.com/bachittle/ping-go/pinger v0.0.0-20201128203717-eb7e3f296fba
	github.com/bachittle/ping-go/utils v0.0.0-20201128215519-db656fffa8fe
)

replace (
	github.com/bachittle/ping-go/pinger => ./pinger
	github.com/bachittle/ping-go/utils => ./utils
)
