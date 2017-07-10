package domain

type baseProvider struct {
	key         string
	name        string
	description string

	metas   []ActionMeta
	actions map[string]Action
}

func NewProvider(key, name, description string) baseProvider {
	b := baseProvider{
		key:         key,
		name:        name,
		description: description,
	}

	Registry.Register(b)

	return b
}

func (b baseProvider) Meta() ProviderMeta {
	return ProviderMeta{b.name, b.description, b.key}
}

func (b baseProvider) ActionList() []ActionMeta {
	return b.metas
}

func (b baseProvider) Action(key string) Action {
	return b.actions[key]
}

func (b baseProvider) Add(a Action) {
	b.actions[a.Meta().Key] = a
	b.metas = append(b.metas, a.Meta())
}
