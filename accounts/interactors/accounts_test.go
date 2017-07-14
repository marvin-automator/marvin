package interactors

import (
	"errors"
	"github.com/bigblind/marvin/accounts"
	"github.com/bigblind/marvin/accounts/domain"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateAccount(t *testing.T) {
	ma := accounts.NewMockAccountStore()
	ma.On("SaveAccount", mock.AnythingOfType("Account")).Return(nil)

	ca := CreateAccount{ma}
	a, err := ca.Execute("foo@example.com", "foo")
	require.NoError(t, err)

	require.Equal(t, "foo@example.com", a.Email)
	ma.AssertExpectations(t)
}

func TestDeleteAccountByID(t *testing.T) {
	ma := accounts.NewMockAccountStore()
	expectedError := errors.New("this was expected")
	ma.On("DeleteAccount", "042").Return(expectedError)

	action := DeleteAccount{ma}
	err := action.ByID("042")
	require.EqualError(t, err, expectedError.Error())
}

func TestDeleteAccountByEmail(t *testing.T) {
	ma := accounts.NewMockAccountStore()
	act := domain.Account{"042", "test@example.com", []byte("nothashed")}
	expectedError := errors.New("this was expected")
	ma.On("GetAccountByEmail", "test@example.com").Return(act, nil)
	ma.On("DeleteAccount", "042").Return(expectedError)

	action := DeleteAccount{ma}
	err := action.ByEmail("test@example.com")
	require.EqualError(t, err, expectedError.Error())
}
