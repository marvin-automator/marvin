package domain

// BasicProvider is a basic implementation of Provider
type BasicProvider struct {
	key         string
	name        string
	description string

	groups           []Group
	defaultGroup     *BasicGroup
	actions          map[string]BaseAction
	globalConfigType interface{}

	// 3rd-party identities
	requiredIdentityProtocol IdentityProtocol
	authorizationEndpoint    string
	tokenEndpoint            string
	requestTokenEndpoint     string
}

// NewProvider creates a new BasicProvider
func NewProvider(key, name, description string) *BasicProvider {
	b := BasicProvider{
		key:         key,
		name:        name,
		description: description,
		actions:     map[string]BaseAction{},
	}

	Registry.Register(&b)

	return &b
}

// SetIdentityParameters sets ProviderMeta values for getting 3rd-party identities
func (b *BasicProvider) SetIdentityParameters(authorizationURL, tokenURL, requestTokenURL string) {
	b.authorizationEndpoint = authorizationURL
	b.tokenEndpoint = tokenURL
	b.requestTokenEndpoint = requestTokenURL
}

// Meta returns metadata about the action provider.
func (b *BasicProvider) Meta() ProviderMeta {
	return ProviderMeta{
		Name:        b.name,
		Description: b.description,
		Key:         b.key,

		ReequiresIdentityProtocol: b.requiredIdentityProtocol,
		AuthorizationEndpoint:     b.authorizationEndpoint,
		TokenEndpoint:             b.tokenEndpoint,
		RequestTokenEndpoint:      b.requestTokenEndpoint,
	}
}

// Groups returns all the groups of this provider.
func (b *BasicProvider) Groups() []Group {
	return b.groups
}

// Action returns the action with the given key
func (b *BasicProvider) Action(key string) BaseAction {
	return b.actions[key]
}

// Add adds an action to the provider
func (b *BasicProvider) Add(a BaseAction) {
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

// NewGroup returns a Group instance that is tied to this provider
func (b *BasicProvider) NewGroup(name string) *BasicGroup {
	g := &BasicGroup{b, name, make([]ActionMeta, 0)}
	b.groups = append(b.groups, g)
	return g
}

// GlobalConfigType returns a struct of the shape that'll hold global configuration
// options for this provider, like API keys. Returns nil if no global configuration
// is necessary.
func (b *BasicProvider) GlobalConfigType() interface{} {
	return b.globalConfigType
}

// Set the global configuration struct returned by GlobalConfigurationType. This is nil by default,
// indicating that no global configuration is necessary.
func (b *BasicProvider) SetGlobalConfigType(c interface{}) {
	b.globalConfigType = c
}

// BasicGroup is a basic implementation of a group.
// Don't create one directly, call BasicProvider.NewGroup instead.
type BasicGroup struct {
	provider *BasicProvider
	name     string
	metas    []ActionMeta
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
func (b *BasicGroup) Add(a BaseAction) {
	b.metas = append(b.metas, a.Meta())
	b.provider.actions[a.Meta().Key] = a
}
