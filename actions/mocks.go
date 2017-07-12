package actions

import (
	"github.com/bigblind/marvin/actions/domain"
	"github.com/stretchr/testify/mock"
)

// MockProvider is a Provider implementation used in tests
type MockProvider struct {
	mock.Mock
}

func NewMockProvider() MockProvider {
	return MockProvider{mock.Mock{}}
}

func (m MockProvider) Meta() domain.ProviderMeta {
	args := m.Called()
	return domain.ProviderMeta{args.String(0), args.String(1), args.String(2)}
}

func (m MockProvider) Groups() []domain.Group {
	args := m.Called()
	return args.Get(0).([]domain.Group)
}

func (m MockProvider) Action(key string) domain.Action {
	args := m.Called(key)
	return args.Get(0).(domain.Action)
}

// MockRegistry is a Registry implementation used in tests
type MockRegistry struct {
	mock.Mock
}

func NewMockRegistr() MockRegistry {
	return MockRegistry{mock.Mock{}}
}

func (m MockRegistry) Register(p domain.ActionProvider) {
	m.Called(p)
}

func (m MockRegistry) Providers() []domain.ProviderMeta {
	args := m.Called()
	return args.Get(0).([]domain.ProviderMeta)
}

func (m MockRegistry) Provider(key string) domain.ActionProvider {
	args := m.Called(key)
	return args.Get(0).(domain.ActionProvider)
}

// MockGroup is a Group implementation used in tests
type MockGroup struct {
	mock.Mock
}

func NewMockGroup() MockGroup {
	return MockGroup{mock.Mock{}}
}

func (m MockGroup) Actions() []domain.ActionMeta {
	args := m.Called()
	return args.Get(0).([]domain.ActionMeta)
}

func (m MockGroup) Name() string {
	args := m.Called()
	return args.String(0)
}
