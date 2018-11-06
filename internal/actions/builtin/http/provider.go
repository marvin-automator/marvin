package http

import (
	"github.com/marvin-automator/marvin/actions"
)

func init() {
	p := actions.Registry.AddProvider("http", "Send and Receive HTTP(S) requests.", []byte{})
	g := p.AddGroup("request", "", []byte{})
	g.AddAction("send", "Send an HTTP Request", []byte{}, makeRequest)
}
