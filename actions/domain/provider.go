package domain

// Internal provider implementation
type baseProvider struct {
	key         string
	name        string
	description string

	metas   []ActionMeta
	actions map[string]Action
}

// Create a new ActionProvider
func NewProvider(key, name, description string) baseProvider {
	b := baseProvider{
		key:         key,
		name:        name,
		description: description,
	}

	Registry.Register(b)

	return b
}

// Meta returns metadata about the action provider.
func (b baseProvider) Meta() ProviderMeta {
	return ProviderMeta{b.name, b.description, b.key}
}

// ActionList returns a list of ActionMeta instances, describing the available
func (b baseProvider) ActionList() []ActionMeta {
	return b.metas
}

// Action returns the action with the given key
func (b baseProvider) Action(key string) Action {
	return b.actions[key]
}

// Add adds an action to the provider
func (b baseProvider) Add(a Action) {
	b.actions[a.Meta().Key] = a
	b.metas = append(b.metas, a.Meta())
}
