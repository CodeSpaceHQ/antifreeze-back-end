package mux

import (
	"github.com/NilsG-S/antifreeze-back-end/common"
)

// maps allow concurrent access so this is fine
type Mux struct {
	routes map[string]func(common.Message, interface{})
}

func New() *Mux {
	return &Mux{
		routes: make(map[string]func(common.Message, interface{})),
	}
}

func (m *Mux) AddRoute(route string, handler func(common.Message, interface{})) {
	m.routes[route] = handler
}

// TODO: this needs to handle operation as well as route. combine the two?
func (m *Mux) Handle(mes common.Message, context interface{}) {
	handler, ok := m.routes[mes.GetSub()]

	if !ok {
		// TODO: avoid swallowing the error?
		return
	}

	handler(mes, context)
}
