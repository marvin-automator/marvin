package interactors

import (
	actiondomain "github.com/marvin-automator/marvin/actions/domain"
	appdomain "github.com/marvin-automator/marvin/app/domain"
	"github.com/marvin-automator/marvin/identityproviders/domain"
	"net/http"
)

// IdentityProvider provides identities for actions.
type IdentityProvider struct {
	Provider      actiondomain.ActionProvider
	IdentityStore domain.IdentityStore
	Logger        appdomain.Logger
}

func (i *IdentityProvider) GetHTTPClient(clientID, clientSecret, account, identity string) *http.Client {
	prot := GetProtocol(IdentityProtocol(i.Provider.Meta().ReequiresIdentityProtocol))
	conf := i.config(clientID, clientSecret)
	id, err := i.IdentityStore.GetIdentity(account, i.Provider.Meta().Key, identity)
	if err != nil {
		i.Logger.Error(err)
		return http.DefaultClient
	}
	cl, err := prot.GetHTTPClient(id.Token, conf)
	if err != nil {
		i.Logger.Error(err)
		return http.DefaultClient
	}
	return cl
}

func (i *IdentityProvider) config(clientID, clientSecret string) ProtocolConfig {
	return ProtocolConfig{
		Endpoint: Endpoint{
			AuthorizationURL: i.Provider.Meta().AuthorizationEndpoint,
			AccessTokenURL:   i.Provider.Meta().TokenEndpoint,
			RequestTokenURL:  i.Provider.Meta().RequestTokenEndpoint,
		},
		Consumer:    clientID,
		Secret:      clientSecret,
		Scopes:      i.Provider.Meta().Scopes,
		CallbackURL: MakeCallbackURL(i.Provider),
	}
}

func MakeCallbackURL(p actiondomain.ActionProvider) string {
	return "/auth/callbacks/" + p.Meta().Name
}
