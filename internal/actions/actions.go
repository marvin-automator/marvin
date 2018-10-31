package actions

import (
	"context"
	"fmt"
	"github.com/marvin-automator/marvin/actions"
	"reflect"
)

type action struct {
	info    actions.Info
	runFunc reflect.Value
}

func (a *action) Info() actions.Info {
	return a.info
}

func (a *action) Run(input interface{}, ctx context.Context) (interface{}, error) {
	retvals := a.runFunc.Call([]reflect.Value{reflect.ValueOf(input)})
	res := retvals[0].Interface()
	err := retvals[1].Interface().(error)
	return res, err
}

func (a *action) validate() {
	name := a.Info().Name

	if a.runFunc.Kind() != reflect.Func {
		panic(fmt.Sprintf("Action %v did not receive a function as runFunc", name))
	}

	ft := a.runFunc.Type()
	ctx := context.Background()
	if !(ft.NumIn() == 2 &&
		ft.In(0).Kind() == reflect.Struct &&
		reflect.TypeOf(ctx).AssignableTo(ft.In(1))) {
		panic(fmt.Sprintf("Action %v should have a function that takes 2 arguments. The first is a struct type that you define, the second is a context.Context", name))
	}

	if a.info.IsTrigger {
		a.validateTrigger(ft)
		a.info.OutputType = ft.Out(0).Elem()
	} else {
		a.validateAction(ft)
		a.info.OutputType = ft.Out(0)
	}

	a.info.InputType = ft.In(0)
}

func (a *action) validateAction(ft reflect.Type) {
	name := a.Info().Name

	var e *error
	if !(ft.NumOut() == 2 && ft.Out(0).Kind() == reflect.Struct && ft.Out(1).Implements(reflect.TypeOf(e).Elem())) {
		panic(fmt.Sprintf("Action %v should have a function that returns 2 values, one of a struct type that you define, and an error.", name))
	}
}

func (a *action) validateTrigger(ft reflect.Type) {
	name := a.Info().Name

	var e *error
	if !(ft.NumOut() == 2 && ft.Out(0).Kind() == reflect.Chan && ft.Out(1).Implements(reflect.TypeOf(e).Elem()) &&
		ft.Out(0).ChanDir()&reflect.RecvDir == reflect.RecvDir) {
		panic(fmt.Sprintf("Trigger %v should have a function that returns 2 values:\n- A readable channel containing structs representing events\n - an error in case anything went wrong.", name))
	}
}

type Group struct {
	actions.BaseInfo
	actions map[string]actions.Action
}

type Provider struct {
	actions.BaseInfo
	groups map[string]*Group
}

func (p *Provider) AddGroup(name, description string, svgIcon []byte) actions.Group {
	g := &Group{actions.BaseInfo{name, description, svgIcon}, make(map[string]actions.Action)}
	p.groups[name] = g
	return g
}

func (p *Provider) Groups() []actions.Group {
	res := make([]actions.Group, 0, len(p.groups))
	for _, g := range p.groups {
		res = append(res, g)
	}

	return res
}

func (g *Group) addAction(name, description string, svgIcon []byte, runFunc interface{}, trigger bool) {
	info := actions.Info{
		BaseInfo:  actions.BaseInfo{name, description, svgIcon},
		IsTrigger: trigger,
	}

	a := &action{
		info:    info,
		runFunc: reflect.ValueOf(runFunc),
	}

	a.validate()

	g.actions[name] = a
}

func (g *Group) AddAction(name, description string, svgIcon []byte, runFunc interface{}) {
	g.addAction(name, description, svgIcon, runFunc, false)
}

func (g *Group) AddManualTrigger(name, description string, svgIcon []byte, runFunc interface{}) {
	g.addAction(name, description, svgIcon, runFunc, true)
}

func (g *Group) Actions() []actions.Action {
	res := make([]actions.Action, 0, len(g.actions))
	for _, a := range g.actions {
		res = append(res, a)
	}

	return res
}

type ProviderRegistry struct {
	providers map[string]*Provider
}

func NewRegistry() *ProviderRegistry {
	return &ProviderRegistry{make(map[string]*Provider)}
}

func (r *ProviderRegistry) AddProvider(name, description string, svgIcon []byte) actions.Provider {
	p := &Provider{
		BaseInfo: actions.BaseInfo{name, description, svgIcon},
		groups:   make(map[string]*Group),
	}

	r.providers[p.Name] = p

	return p
}

func (r *ProviderRegistry) Providers() []actions.Provider {
	ps := make([]actions.Provider, 0, len(r.providers))
	for _, p := range r.providers {
		ps = append(ps, p)
	}

	return ps
}

func (r *ProviderRegistry) GetAction(provider, group, action string) (actions.Action, error) {
	if p, ok := r.providers[provider]; ok {
		if g, ok := p.groups[group]; ok {
			if a, ok := g.actions[action]; ok {
				return a, nil
			}
			return nil, fmt.Errorf("Group %v->%v has no action %v", provider, group, action)
		}
		return nil, fmt.Errorf("Provider %v has no group %v", provider, group)
	}
	return nil, fmt.Errorf("no provider: %v", provider)
}

func init() {
	actions.Registry = NewRegistry()
}
