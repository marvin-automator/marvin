package handlers

import (
	"github.com/marvin-automator/marvin/accounts"
	"github.com/marvin-automator/marvin/actions/interactors"
	actionstorage "github.com/marvin-automator/marvin/actions/storage"
	"github.com/marvin-automator/marvin/handlers"
)

// The AccountChores handler responds with the current account's chores in JSON.
func AccountChores(c handlers.Context) error {
	s := c.Store()
	cs := actionstorage.NewChoreStore(s)
	in := interactors.GetChores{cs}
	chores, err := in.ForAccount(accounts.CurrentAccount(c).ID)

	if err != nil {
		return c.Error(500, err)
	}

	return c.Render(200, c.Renderer().JSON(chores))
}
