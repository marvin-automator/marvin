package accounts

import (
	"github.com/marvin-automator/marvin/handlers"

	"github.com/marvin-automator/marvin/accounts/interactors"
)

// CurrentAccount returns the authenticated account for the request, associated with the given context.
func CurrentAccount(c handlers.Context) interactors.Account {
	if ac, ok := c.Value("account").(interactors.Account); ok {
		return ac
	}
	return interactors.Account{}
}
