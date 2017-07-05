package interactors

import (
	"github.com/bigblind/marvin/accounts/domain"
)

// The Account type that the actions return to handlers.
type Account struct {
	ID    string
	Email string
}

type CreateAccount struct {
	AccountStore domain.AccountStore
}

func (c CreateAccount) Execute(email, password string) (Account, error) {
	act, err := domain.NewAccount(email, password)
	if err != nil {
		return Account{}, err
	}

	c.AccountStore.SaveAccount(act)

	return Account{act.ID, act.Email}, nil
}

type ICreateAccount interface {
	Execute(email, password string) (Account, error)
}

type DeleteAccount struct {
	AccountStore domain.AccountStore
}

func (d DeleteAccount) ByID(aid string) error {
	return d.AccountStore.DeleteAccount(aid)
}

func (d DeleteAccount) ByEmail(email string) error {
	act, err := d.AccountStore.GetAccountByEmail(email)
	if err != nil {
		return err
	}
	return d.AccountStore.DeleteAccount(act.ID)
}
