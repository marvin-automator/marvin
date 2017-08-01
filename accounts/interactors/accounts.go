package interactors

import (
	"github.com/marvin-automator/marvin/accounts/domain"
)

// The Account type that the actions return to handlers.
type Account struct {
	ID    string
	Email string
}

// CreateAccount is an interactor to create an account
type CreateAccount struct {
	AccountStore domain.AccountStore
}

// Execute actually executes the interactor
func (c CreateAccount) Execute(email, password string) (Account, error) {
	act, err := domain.NewAccount(email, password)
	if err != nil {
		return Account{}, err
	}

	err = c.AccountStore.SaveAccount(act)

	return Account{act.ID, act.Email}, err
}

// ICreateAccount is an interface implemented by CreateAccount, so that it can be mocked
type ICreateAccount interface {
	Execute(email, password string) (Account, error)
}

// DeleteAccount is an interactor to delete accounts
type DeleteAccount struct {
	AccountStore domain.AccountStore
}

// ByID executes the interactor, deleting the account with the given ID
func (d DeleteAccount) ByID(aid string) error {
	return d.AccountStore.DeleteAccount(aid)
}

// ByEmail executes the interactor, deleting the account with the given email.
func (d DeleteAccount) ByEmail(email string) error {
	act, err := d.AccountStore.GetAccountByEmail(email)
	if err != nil {
		return err
	}
	return d.AccountStore.DeleteAccount(act.ID)
}

// IDeleteAccount is an interface, implemented by DeleteAccount, so that it can be mocked.
type IDeleteAccount interface {
	ByID(aid string) error
	ByEmail(email string) error
}

// GetAccount is an interactor for getting accounts
type GetAccount struct {
	accountStore domain.AccountStore
}

// ByID returns the account with the given ID.
func (g GetAccount) ByID(aid string) (Account, error) {
	a, err := g.accountStore.GetAccountByID(aid)
	if err != nil {
		return Account{}, err
	}
	return Account{a.ID, a.Email}, nil
}