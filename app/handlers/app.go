package handlers

import "github.com/marvin-automator/marvin/handlers"

// AppPage renders the page that contains the main JavaScript app.
func AppPage(c handlers.Context) error {
	c.Logger().Debug("Rendering JS App view")
	return c.Render(200, c.Renderer().HTML("jsapp.html)"))
}