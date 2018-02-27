// Consider moving this file to a `common` directory, since it'll be used by `rest` and `ws`
package ws

import (
	"time"
)

const (
	add    int = 1
	remove int = 2
	update int = 3
)

// support multiple message types
type message interface {
	getSub() string
	getOp() int
}

type temp struct {
	sub      string
	op       int
	deviceId int
	temp     int
	time     time.Time
}

func (v temp) getSub() string { return v.sub }
func (v temp) getOp() int     { return v.op }
