package handlers

import (
	"github.com/bigblind/marvin/handlers"
	actionstorage "github.com/bigblind/marvin/actions/storage"
	"github.com/bigblind/marvin/actions/interactors"
	"github.com/bigblind/marvin/accounts"
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