// Authored and revised by YOC team, 2017-2018
// License placeholder #1

//go:generate go-bindata -nometadata -o assets.go -pkg tracers -ignore ((tracers)|(assets)).go ./...
//go:generate gofmt -s -w assets.go

// Package tracers contains the actual JavaScript tracer assets.
package tracers
