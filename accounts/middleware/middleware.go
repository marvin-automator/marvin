package middleware

import (
	"github.com/bigblind/marvin/accounts/domain"
	"github.com/bigblind/marvin/accounts/interactors"
	accountsstorage "github.com/bigblind/marvin/accounts/storage"
	configstorage "github.com/bigblind/marvin/config/storage"
	"github.com/bigblind/marvin/handlers"
	"github.com/bigblind/marvin/storage"
	"github.com/gobuffalo/buffalo"
	"github.com/pkg/errors"
)

// The key in the session where the uid of the currently logged-in user is stored.
var uidKey = "login_uid"

// Error used to check if the user needs to be redirected to the login page
var errNeedsLogin = errors.New("errNeedsLogin to login")

// Middleware checks whether the user is logged in, and redirects them to login if necessary.
// It stores the current account in the context.
//todo: make this more testable by passing in some kind of factories for account and config stores.
func Middleware(next buffalo.Handler) buffalo.Handler {
	var h handlers.Handler
	h = func(c handlers.Context) error {
		var account interactors.Account

		// Get the id of the currently logged-in account, if any.
		uid := c.Session().Get(uidKey)
		c.Logger().Debugf("UID in session: %v", uid)

		err := c.WithReadableStore(func(s storage.Store) error {
			// Set up the necessary stores and interactor
			as := accountsstorage.NewAccountStore(s)
			cs := configstorage.NewConfigStore(s)
			i := interactors.Login{as, cs}

			// Check whether we're in accounts-enabled mode.
			req, err := i.IsRequired()
			if err != nil {
				return err
			}

			// If accounts are not enabled, use the default account
			if !req {
				c.Logger().Debug("Login not required")
				account, err = i.GetDefaultAccount()
				if err != nil {
					return err
				}
			} else {
				// If there's no user id in our session, make them log in
				if uid == nil {
					return errNeedsLogin
				}

				// Try to get the user with the id in the session
				account, err = i.GetAccountByID(uid.(string))

				// If there's no user with this ID...
				if err == domain.ErrAccountNotFound {
					// ... Make them log in again
					c.Logger().Debug("No user with ID %v", uid)
					return errNeedsLogin
					//Any other error should be returned as normal
				} else if err != nil {
					return err
				}
			}

			// Save the account in the context and session.
			c.Logger().Debug("Account found, store it in the session")
			c.Set("account", account)
			c.Session().Set(uidKey, account.ID)
			return nil
		})

		if err == errNeedsLogin {
			return c.Redirect(302, "/login")
		}
		return next(c)
	}

	return h.ToBuffalo()
}
