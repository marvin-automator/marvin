package oauth2

import (
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal/oauth2"
	go_oauth "golang.org/x/oauth2"
)

func Requirement(endpoint go_oauth.Endpoint, helpTemplate string, getAccount AccountGetter) actions.Requirement {
	return &oauth2.OAuth{
		Endpoint: endpoint,
		ConfigHelpTemplate: helpTemplate,
		GatAccount: func(token *go_oauth.Token) (oauth2.Account, error) {
			acc, err := getAccount(token)

			return oauth2.Account{
				Token: acc.Token,
				Name: acc.Name,
				Id: acc.Id,
				ImageURL: acc.ImageURL,
			}, err
		},
	}
}