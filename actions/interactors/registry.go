package interactors

import (
	"github.com/marvin-automator/marvin/actions/domain"
)

// The Registry interactor gives access to the available actions.
type Registry struct {
	r domain.ProviderRegistry
}

// Group represents a group of actions.
type Group struct {
	Name     string              `json:"name"`
	Provider string              `json:"provider"`
	actions  []domain.ActionMeta `json:"actions"`
}

func (g Group) filterActions(isTrigger bool) []domain.ActionMeta {
	res := []domain.ActionMeta{}
	for _, act := range g.actions {
		if act.IsTrigger == isTrigger {
			res = append(res, act)
		}
	}

	return res
}

// Actions returns all the actions, excluding triggers, in the group.
func (g Group) Actions() []domain.ActionMeta { return g.filterActions(false) }

// Triggers returns all the triggers in the group.
func (g Group) Triggers() []domain.ActionMeta { return g.filterActions(true) }

type Provider struct {
	Key    string  `json:"key"`
	Name   string  `json:"name"`
	Groups []Group `json:"groups"`
}

// NewRegistryInteractor returns a new instance of the Registry interactors
func NewRegistryInteractor() Registry {
	return Registry{domain.Registry}
}

// GetActionGroups returns a list of available groups
func (r Registry) GetActionGroups() []Group {
	gs := make([]Group, 0)
	for _, pm := range r.r.Providers() {
		p := r.r.Provider(pm.Key)
		for _, g := range p.Groups() {
			gs = append(gs, Group{g.Name(), pm.Key, g.Actions()})
		}
	}

	return gs
}

func (r Registry) GetProviders() []Provider {
	ps := make([]Provider, 0)
	for _, pm := range r.r.Providers() {
		p := Provider{pm.Key, pm.Name, []Group{}}
		for _, g := range r.r.Provider(pm.Key).Groups() {
			p.Groups = append(p.Groups, Group{g.Name(), pm.Key, g.Actions()})
		}
		ps = append(ps, p)
	}

	return ps
}
