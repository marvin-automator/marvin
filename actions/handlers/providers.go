package handlers

import (
	"github.com/marvin-automator/marvin/actions/interactors"
	"github.com/marvin-automator/marvin/handlers"
)

// ActionGroups responds with the available groups with their actions, in JSON.
func ActionGroups(c handlers.Context) error {
	int := interactors.NewRegistryInteractor()
	return c.Render(200, c.Renderer().JSON(int.GetActionGroups()))
}

func SpriteSheet(c handlers.Context) error {
	sheet, err := interactors.NewRegistryInteractor().GetSpriteSheet()
	if err != nil {
		return c.Error(500, err)
	}

	c.Response().Header().Set("Content-Type", "image/svg+xml")
	c.Response().WriteHeader(200)
	_, err = c.Response().Write(sheet)
	return err
}
