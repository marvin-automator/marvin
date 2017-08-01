package handlers

import (
	"github.com/marvin-automator/marvin/handlers"
	"github.com/marvin-automator/marvin/actions/interactors"
)

// The ActionGroups handler responds with the available groups with their actions, in JSON.
func ActionGroups(c handlers.Context) error {
	int := interactors.NewRegistryInteractor()
	return c.Render(200, c.Renderer().JSON(int.GetActionGroups()))
}