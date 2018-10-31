package http

import (
	"github.com/marvin-automator/marvin/actions"
)

func init() {
	p := actions.Registry.AddProvider("HTTP", "Send and Receive HTTP(S) requests.", []byte{})
	g := p.AddGroup("default", "", []byte{})
	g.AddAction("send request", "Send an HTTP Request", []byte{}, makeRequest)
}
