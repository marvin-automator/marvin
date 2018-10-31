package actions

import (
	"context"
	"reflect"
)

type BaseInfo struct {
	Name        string
	Description string
	SVGIcon     []byte
}

func (i BaseInfo) Info() BaseInfo {
	return i
}

type Info struct {
	BaseInfo
	InputType  reflect.Type
	OutputType reflect.Type
	IsTrugger  bool
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
	Info() BaseInfo
}

type ProviderRegistry interface {
	AddProvider(name, description string, svgIcon []byte) Provider
	GetAction(provider, group, action string) (Action, error)
	Providers() []Provider
}

var Registry ProviderRegistry
