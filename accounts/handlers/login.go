package handlers

import (
	"github.com/bigblind/marvin/handlers"
	"github.com/gobuffalo/buffalo/render"
)


// The key in the session under which the user ID is stored
var uidKey = "login_uid"

func LoginPage(c handlers.Context) error {
	email, _ := c.Session().Get("login_email").(string)

	return c.BareRenderer().HTML("accounts/login.html").Render(c.Response(), render.Data{
		"login_email": email,
	})
}
