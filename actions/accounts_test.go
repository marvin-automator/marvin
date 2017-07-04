package actions

import (
	"testing"
	"github.com/bigblind/marvin/mocks"
	"github.com/bigblind/marvin/domain"
	"github.com/stretchr/testify/require"
	"errors"
	"github.com/stretchr/testify/mock"
)

func TestSuccessfulLogin(t *testing.T) {
	acc, err := domain.NewAccount("test@example.com", "pwd")
	require.NoError(t, err)
	exp := Account{acc.ID,"test@example.com"}
	ma := mocks.NewMockAccountStore()
	ma.On("GetAccountByEmail", "test@example.com").Return(acc, nil)
	mc := mocks.MockConfigStore{domain.DefaultConfig, nil}

	login := Login{ma,mc}
	res, err := login.Execute("test@example.com", "pwd")
	require.Equal(t, exp, res)
}

func TestWrongPassword(t *testing.T) {
	acc, err := domain.NewAccount("test@example.com", "pwd")
	ma := mocks.NewMockAccountStore()
	ma.On("GetAccountByEmail", "test@example.com").Return(acc, nil)
	mc := mocks.MockConfigStore{domain.DefaultConfig, nil}
	require.NoError(t, err)

	login := Login{ma,mc}
	_, err = login.Execute("test@example.com", "incorrect")
	require.EqualError(t, err, ErrLoginFailed.Error())
}

func TestAccountNotFoundReturnsFailedLogin(t *testing.T) {
	ma := mocks.NewMockAccountStore()
	ma.On("GetAccountByEmail", "test@example.com").Return(domain.Account{}, domain.ErrAccountNotFound)
	mc := mocks.MockConfigStore{domain.DefaultConfig, nil}

	login := Login{ma,mc}
	_, err := login.Execute("test@example.com", "pwd")
	require.EqualError(t, err, ErrLoginFailed.Error(), "We need to return the same error as if the password was wrong.")
}

func TestConfigStoreError(t *testing.T) {
	ma := mocks.NewMockAccountStore()
	configError := errors.New("something went wrong")
	mc := mocks.MockConfigStore{domain.DefaultConfig, configError}

	login := Login{ma,mc}
	_, err := login.Execute("test@example.com", "pwd")
	require.EqualError(t, err, configError.Error())
}

func TestAccountStoreError(t *testing.T) {
	accountError := errors.New("something went wrong")
	ma := mocks.NewMockAccountStore()
	ma.On("GetAccountByEmail", "test@example.com").Return(domain.Account{}, accountError)
	mc := mocks.MockConfigStore{domain.DefaultConfig, nil}

	login := Login{ma,mc}
	_, err := login.Execute("test@example.com", "pwd")
	require.EqualError(t, err, accountError.Error())
}

func TestAccountsDisabled(t *testing.T) {
	ma := mocks.NewMockAccountStore()
	mc := mocks.MockConfigStore{domain.Config{AccountsEnabled: false}, nil}

	login := Login{ma,mc}
	_, err := login.Execute("test@example.com", "pwd")
	require.EqualError(t, err, ErrAccountsDisabled.Error())
}

func TestCreateAccount(t *testing.T) {
	ma := mocks.NewMockAccountStore()
	ma.On("SaveAccount", mock.AnythingOfType("Account")).Return(nil)

	ca := CreateAccount{ma}
	a, err := ca.Execute("foo@example.com", "foo")
	require.NoError(t, err)

	require.Equal(t, "foo@example.com", a.Email)
	ma.AssertExpectations(t)
}