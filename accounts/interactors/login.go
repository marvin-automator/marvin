package interactors

import (
	"errors"

	"github.com/marvin-automator/marvin/accounts/domain"
	configdomain "github.com/marvin-automator/marvin/config/domain"
)

// ErrAccountsDisabled is returned when trying to log in when accounts are disabled
var ErrAccountsDisabled = errors.New("can not log in when accounts are disabled")

// ErrLoginFailed is returned when the password doesn't match the account's password, or when the account isn't found
var ErrLoginFailed = errors.New("incorrect email or password")

// Login is an interactor for handling logging in
type Login struct {
	AccountStore domain.AccountStore
	ConfigStore  configdomain.ConfigStore
}

// Execute executed the interactor, returning the account that matches the given email address and password, if any.
// If accounts are disabled, returns ErrAccountsDisabled. If the email/password combination is incorrect,
// this function returns ErrLoginFailed
func (l Login) Execute(email, password string) (Account, error) {
	config, err := l.ConfigStore.GetConfig()
	if err != nil {
		return Account{}, err
	}

	if !config.AccountsEnabled {
		return Account{}, ErrAccountsDisabled
	}

	act, err := l.AccountStore.GetAccountByEmail(email)

	if err == domain.ErrAccountNotFound {
		return Account{}, ErrLoginFailed
	}

	if err != nil {
		return Account{}, err
	}

	result, err := act.CheckPassword(password)
	if err != nil {
		return Account{}, err
	}

	if !result {
		return Account{}, ErrLoginFailed
	}

	return Account{act.ID, act.Email}, nil
}

// IsRequired returns whether login is required by the current configuration
func (l Login) IsRequired() (bool, error) {
	c, err := l.ConfigStore.GetConfig()
	if err != nil {
		return true, err
	}
	return c.AccountsEnabled, nil
}

// GetAccountByID returns the account with the given ID.
func (l Login) GetAccountByID(aid string) (a Account, err error) {
	da, err := l.AccountStore.GetAccountByID(aid)
	if err == nil {
		a = Account{da.ID, da.Email}
	}
	return
}

// GetDefaultAccount returns the default account
func (l Login) GetDefaultAccount() (a Account, err error) {
	da, err := l.GetDefaultAccount()
	if err == nil {
		a = Account{da.ID, da.Email}
	}
	return
}

// ILogin is the interface implemented by Login, so that it can be mocked.
type ILogin interface {
	Execute(email, password string) (bool, error)
	IsRequired() (bool, error)
	GetAccountByID(aid string) (Account, error)
	GetDefaultAccount() (Account, error)
}
