package domain

import "sync"

var Registry = new(registry)

type registry struct {
	providers map[string]ActionProvider
	mut sync.Mutex
}

func (r *registry) Register(p ActionProvider) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.providers[p.Meta().Key] = p
}

// Returns a slice of available providers
func (r *registry) GetProviders() []ProviderMeta {
	l := make([]ProviderMeta, len(r.providers))
	for _, p := range r.providers {
		l = append(l, p.Meta())
	}

	return l
}

func(r *registry) Provider(key string) ActionProvider {
	return r.providers[key]
}