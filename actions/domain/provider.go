package domain

// BasicProvider is a basic implementation of Provider
type BasicProvider struct {
	key         string
	name        string
	description string

	groups  []Group
	defaultGroup *BasicGroup
	actions map[string]Action
}

// Create a new ActionProvider
func NewProvider(key, name, description string) *BasicProvider {
	b := BasicProvider{
		key:         key,
		name:        name,
		description: description,
	}

	Registry.Register(&b)

	return &b
}

// Meta returns metadata about the action provider.
func (b *BasicProvider) Meta() ProviderMeta {
	return ProviderMeta{b.name, b.description, b.key}
}

func (b *BasicProvider) Groups() []Group {
	return b.groups
}

// Action returns the action with the given key
func (b *BasicProvider) Action(key string) Action {
	return b.actions[key]
}

// Add adds an action to the provider
func (b *BasicProvider) Add(a Action) {
	g := b.getOrCreateDefaultGroup()
	g.Add(a)
}

// getOrCreateDefaultGroup returns the default group of this provider, creating it if it doesn't exist.
func (b *BasicProvider) getOrCreateDefaultGroup() *BasicGroup {
	if b.defaultGroup == nil {
		b.defaultGroup = b.NewGroup(b.Meta().Name)
	}
	return b.defaultGroup
}

// NewGroup returns a Group instance that is tie
func (b *BasicProvider) NewGroup(name string) *BasicGroup {
	g := &BasicGroup{b, name, make([]ActionMeta, 0)}
	b.groups = append(b.groups, g)
	return g
}

// BasicGroup is a basic implementation of a group.
// Don't create one directly, call BasicProvider.NewGroup instead.
type BasicGroup struct {
	provider *BasicProvider
	name string
	metas []ActionMeta
}

// Name returns a human-readable name for the group.
func (b *BasicGroup) Name() string {
	return b.name
}

// Actions returns ActionMetas describing the actions in this group.
func (b *BasicGroup) Actions() []ActionMeta {
	return b.metas
}

// Add adds an action to this group.
func (b *BasicGroup) Add(a Action) {
	b.metas = append(b.metas, a.Meta())
	b.provider.actions[a.Meta().Key] = a
}

