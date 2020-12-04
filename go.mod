module github.com/bachittle/ping-go

go 1.15

require (
	github.com/bachittle/ping-go/pinger v0.0.0-20201204042230-3e400ffa0019
	github.com/bachittle/ping-go/utils v0.0.0-20201204042230-3e400ffa0019
)

replace (
	github.com/bachittle/ping-go/pinger => ./pinger
	github.com/bachittle/ping-go/utils => ./utils
)
