package actions

import (
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal/oauth2"
	original_oauth2 "golang.org/x/oauth2"
)

// Requirements specify the things an action needs to be able to run.
type Requirement interface {
	Init(p actions.Provider)
	Config() interface{}
	ConfigHelp() string
	SetConfig(interface{}) error
	Fulfilled() bool
}

type Requirements struct {
	OAuth2 *oauth2.OAuth
}

func (r *Requirements) AddOAuth2(p actions.Provider, endpoint original_oauth2.Endpoint, helpTemplate string, getAccount oauth2.AccountGetter) {
	r.OAuth2 = &oauth2.OAuth{
		Endpoint: endpoint,
		ConfigHelpTemplate: helpTemplate,
		Provider: p,
		GatAccount: getAccount,
	}
}







