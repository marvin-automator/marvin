package handlers

import (
	"github.com/bigblind/marvin/handlers"
	"github.com/bigblind/marvin/storage"
	"github.com/bigblind/marvin/actions/interactors"
)

// The ActionGroups handler responds with the available groups with their actions, in JSON.
func ActionGroups(c handlers.Context) error {
	return c.WithReadableStore(func(s storage.Store) error {
		int := interactors.NewRegistryInteractor()
		return c.Render(200, c.Renderer().JSON(int.GetActionGroups()))
	})
}