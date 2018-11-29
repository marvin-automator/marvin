package actions

import (
	"context"
	"golang.org/x/oauth2"
	"reflect"
	"strings"
)

type BaseInfo struct {
	Name        string
	Description string
	SVGIcon     []byte
	Parent		*BaseInfo
}

func (i BaseInfo) Info() BaseInfo {
	return i
}

func (i BaseInfo) Path() string {
	names := []string{i.Name}
	for p := i.Parent; p != nil; p = p.Parent {
		names = append([]string{p.Name}, names...)
	}
	return strings.Join(names, ".")
}

type Info struct {
	BaseInfo
	InputType  reflect.Type
	OutputType reflect.Type
	IsTrigger  bool
}

type Action interface {
	Info() Info
	Run(input interface{}, ctx context.Context) (interface{}, error)
}

type Group interface {
	AddAction(name, description string, svgInput []byte, runFunc interface{})
	AddManualTrigger(name, description string, svgIcon []byte, runFunc interface{})
	Actions() []Action
	Info() BaseInfo
}

type Provider interface {
	AddGroup(name, description string, svgIcon []byte) Group
	Groups() []Group

	SetOAuth2Endpoint(endpoint oauth2.Endpoint)
	Info() BaseInfo
}

type ProviderRegistry interface {
	AddProvider(name, description string, svgIcon []byte) Provider
	GetAction(provider, group, action string) (Action, error)
	Providers() []Provider
}

var Registry ProviderRegistry
