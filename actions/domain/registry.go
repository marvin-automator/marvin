package domain

import "sync"

// Registry maps provider keys to provider instances
var Registry ProviderRegistry = newRegistry()

func newRegistry() *registry {
	r := registry{
		providers: map[string]ActionProvider{},
	}
	return &r
}

type registry struct {
	providers map[string]ActionProvider
	mut       sync.Mutex
}

// ProviderRegistry is the interface implemented by Registry
type ProviderRegistry interface {
	// Register registers a new ActionProvider
	Register(p ActionProvider)
	// Providers returns a list of ProviderMeta instances,
	// describing the available action providers
	Providers() []ProviderMeta
	// Provider returns the ActionProvider with the given key
	Provider(key string) ActionProvider
}

func (r *registry) Register(p ActionProvider) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.providers[p.Meta().Key] = p
}

// Providers returns a slice of available providers
func (r *registry) Providers() []ProviderMeta {
	l := make([]ProviderMeta, 0)
	for _, p := range r.providers {
		l = append(l, p.Meta())
	}

	return l
}

// Provider returns the ActionProvider with the given key
func (r *registry) Provider(key string) ActionProvider {
	return r.providers[key]
}
