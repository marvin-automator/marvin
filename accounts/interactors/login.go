package interactors

import (
	"errors"

	"github.com/bigblind/marvin/accounts/domain"
	configdomain "github.com/bigblind/marvin/config/domain"
)

// Returned when trying to log in when accounts are disabled
var ErrAccountsDisabled = errors.New("Can not log in when accounts are disabled.")

// Returned when the password doesn't match the account's password, or when the account isn't found
var ErrLoginFailed = errors.New("Incorrect email or password")

type Login struct {
	AccountStore domain.AccountStore
	ConfigStore  configdomain.ConfigStore
}

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

func(l Login) IsRequired() (bool, error) {
	c, err := l.ConfigStore.GetConfig();
	if err != nil {
		return true, err
	}
	return c.AccountsEnabled, nil
}

func(l Login) GetAccountByID(aid string) (a Account, err error) {
	da, err := l.AccountStore.GetAccountByID(aid)
	if err == nil {
		a = Account{da.ID,da.Email}
	}
	return
}

func(l Login) GetDefaultAccount() (a Account, err error) {
	da, err := l.GetDefaultAccount()
	if err == nil {
		a = Account{da.ID,da.Email}
	}
	return
}

type ILogin interface {
	Execute(email, password string) (bool, error)
	IsRequired() (bool, error)
	GetAccountByID(aid string) (Account, error)
	GetDefaultAccount() (Account, error)
}
