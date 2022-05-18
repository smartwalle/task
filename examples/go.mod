module github.com/smartwalle/task/examples

require (
	github.com/smartwalle/task v0.0.0
	github.com/smartwalle/queue v0.0.2
)

replace (
	github.com/smartwalle/task => ../
)

go 1.18
