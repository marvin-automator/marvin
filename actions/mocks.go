package actions

import (
	"github.com/marvin-automator/marvin/actions/domain"
	"github.com/stretchr/testify/mock"
)

// MockProvider is a Provider implementation used in tests
type MockProvider struct {
	*mock.Mock
}

// NewMockProvider returns a new MockProvider
func NewMockProvider() *MockProvider {
	mo := mock.Mock{}
	mp := MockProvider{&mo}
	return &mp
}

// Meta is the mocked implementation of Provider.Meta
func (m *MockProvider) Meta() domain.ProviderMeta {
	args := m.Called()
	return domain.ProviderMeta{
		Name: args.String(0),
		Key: args.String(1),
		Description: args.String(2),
	}
}

// Groups is the mocked implementation of Provider.Groups
func (m *MockProvider) Groups() []domain.Group {
	args := m.Called()
	return args.Get(0).([]domain.Group)
}

// Action is the mocked implementation of Provider.Action
func (m *MockProvider) Action(key string) domain.BaseAction {
	args := m.Called(key)
	return args.Get(0).(domain.Action)
}

// GlobalConfigType is the mocked implementation of Provider.GlobalConfigType
func (m *MockProvider) GlobalConfigType() interface{} {
	args := m.Called()
	return args.Get(0)
}

// MockRegistry is a Registry implementation used in tests
type MockRegistry struct {
	*mock.Mock
}

// NewMockRegistry returns a new MockRegistry
func NewMockRegistry() *MockRegistry {
	mo := mock.Mock{}
	mr := MockRegistry{&mo}
	return &mr
}

// Register is the mocked implementation of ProviderRegistry.Register
func (m *MockRegistry) Register(p domain.ActionProvider) {
	m.Called(p)
}

// Providers is the mocked implementation of ProviderRegistry.Providers
func (m *MockRegistry) Providers() []domain.ProviderMeta {
	args := m.Called()
	return args.Get(0).([]domain.ProviderMeta)
}

// Provider is the mocked implementation of ProviderRegistry.Provider
func (m *MockRegistry) Provider(key string) domain.ActionProvider {
	args := m.Called(key)
	return args.Get(0).(domain.ActionProvider)
}

// MockGroup is a Group implementation used in tests
type MockGroup struct {
	*mock.Mock
}

// NewMockGroup returns a new MockGroup
func NewMockGroup() *MockGroup {
	mo := mock.Mock{}
	mg := MockGroup{&mo}
	return &mg
}

// Actions is the mocked implementation of Group.Actions
func (m *MockGroup) Actions() []domain.ActionMeta {
	args := m.Called()
	return args.Get(0).([]domain.ActionMeta)
}

// Name is the mocked implementation of Group.Name
func (m *MockGroup) Name() string {
	args := m.Called()
	return args.String(0)
}
