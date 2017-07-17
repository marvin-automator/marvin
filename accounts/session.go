package accounts

import (
	"github.com/bigblind/marvin/handlers"

	"github.com/bigblind/marvin/accounts/interactors"
)

func CurrentAccount(c handlers.Context) interactors.Account {
	if ac, ok := c.Value("account").(interactors.Account); ok {
		return ac
	}
	return interactors.Account{}
}