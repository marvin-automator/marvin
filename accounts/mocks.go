package accounts

import (
	"github.com/bigblind/marvin/accounts/domain"
	"github.com/stretchr/testify/mock"
)

type MockAccountStore struct {
	mock.Mock
}

func (a MockAccountStore) SaveAccount(account domain.Account) error {
	vals := a.Called(account)
	return vals.Error(0)
}

func NewMockAccountStore() MockAccountStore {
	return MockAccountStore{mock.Mock{}}
}

func (a MockAccountStore) GetAccountByID(aid string) (domain.Account, error) {
	vals := a.Called(aid)
	return vals.Get(0).(domain.Account), vals.Error(1)
}

func (a MockAccountStore) GetAccountByEmail(email string) (domain.Account, error) {
	vals := a.Called(email)
	return vals.Get(0).(domain.Account), vals.Error(1)
}

func (a MockAccountStore) DeleteAccount(aid string) error {
	vals := a.Called(aid)
	return vals.Error(0)
}
