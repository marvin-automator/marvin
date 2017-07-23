package mocks

import (
	"github.com/bigblind/marvin/accounts/domain"
	"github.com/stretchr/testify/mock"
)

// MockAccountStore is a mock implementation of the AccountStore interface
type MockAccountStore struct {
	*mock.Mock
}

// NewMockAccountStore returns a new MockAccountStore
func NewMockAccountStore() *MockAccountStore {
	mo := mock.Mock{}
	ma := MockAccountStore{&mo}
	return &ma
}

// SaveAccount mocks the SaveAccount method of the AccountStore interface
func (a *MockAccountStore) SaveAccount(account domain.Account) error {
	vals := a.Called(account)
	return vals.Error(0)
}

// GetAccountByID mocks the GetAccountByID method of the AccountStore interface.
func (a *MockAccountStore) GetAccountByID(aid string) (domain.Account, error) {
	vals := a.Called(aid)
	return vals.Get(0).(domain.Account), vals.Error(1)
}

// GetAccountByEmail mocks the GetAccountByEmail method of the AccountStore interface.
func (a *MockAccountStore) GetAccountByEmail(email string) (domain.Account, error) {
	vals := a.Called(email)
	return vals.Get(0).(domain.Account), vals.Error(1)
}

// DeleteAccount mocks the DeleteAccount method of the AccountStore interface.
func (a *MockAccountStore) DeleteAccount(aid string) error {
	vals := a.Called(aid)
	return vals.Error(0)
}


// EachAccount mocks the EachAcount method of the AccountStore interface
func(a *MockAccountStore) EachAccount(f func(domain.Account) error) error {
	args := a.Called(f)
	return args.Error(0)
}