package handlers

import (
	"github.com/bigblind/marvin/accounts/interactors"
	accountsstorage "github.com/bigblind/marvin/accounts/storage"
	configstorage "github.com/bigblind/marvin/config/storage"
	"github.com/bigblind/marvin/handlers"
	"github.com/bigblind/marvin/storage"
	"net/url"
)

// The key in the session under which the user ID is stored
var uidKey = "login_uid"

// LoginPage handles the GET request to the login route, rendering the form.
func LoginPage(c handlers.Context) error {
	email, _ := c.Session().Get("login_email").(string)
	c.Session().Delete("login_email")
	c.Set("login_email", email)
	err := c.Param("error")
	c.Set("error", err)
	return c.Render(200, c.BareRenderer().HTML("accounts/login.html"))
}

// ProcessLogin handles POST requests to the login URL, processing form information.
func ProcessLogin(c handlers.Context) error {
	email := c.Request().Form.Get("email")
	password := c.Request().Form.Get("password")

	return c.WithReadableStore(func(s storage.Store) error {
		as := accountsstorage.NewAccountStore(s)
		cs := configstorage.NewConfigStore(s)
		login := interactors.Login{AccountStore: as, ConfigStore: cs}

		act, err := login.Execute(email, password)
		if err == interactors.ErrLoginFailed {
			c.Session().Set("login_email", email)
			err = c.Redirect(302, "/login?error="+url.QueryEscape(err.Error()))
			return err
		} else if err != nil {
			return err
		}

		c.Session().Set(uidKey, act.ID)
		err = c.Redirect(302, "/")
		return err
	})
}
