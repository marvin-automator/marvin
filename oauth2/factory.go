package oauth2

import (
	"context"
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal/oauth2"
	go_oauth "golang.org/x/oauth2"
)

func Requirement(endpoint go_oauth.Endpoint, helpTemplate string, getAccount AccountGetter) actions.Requirement {
	var o *oauth2.OAuth

	o = &oauth2.OAuth{
		Endpoint:           endpoint,
		ConfigHelpTemplate: helpTemplate,
		GatAccount: func(token *go_oauth.Token) (oauth2.Account, error) {
			cl := o.GoConfig([]string{}).Client(context.TODO(), token)
			acc, err := getAccount(cl)

			return oauth2.Account{
				Token:    token,
				Name:     acc.Name,
				Id:       acc.Id,
				ImageURL: acc.ImageURL,
			}, err
		},
	}

	return o
}
