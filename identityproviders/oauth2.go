package identityproviders

import (
	"context"
	"github.com/marvin-automator/marvin/identityproviders/interactors"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"net/http"
	"strings"
	"time"
)

type OAuth2 struct{}

func init() {
	interactors.RegisterImplementation(interactors.OAuth2, OAuth2{})
}

// AuthorizationURL returns the URL the user should be redirected to for granting authorization.
func (o OAuth2) AuthorizationURL(conf interactors.ProtocolConfig, state string, c context.Context) (string, error) {
	return makeOAuth2Config(conf).AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce), nil
}

// GetToken processes the http request from the callback URL, and returns a token
func (o OAuth2) GetToken(req http.Request, state string, conf interactors.ProtocolConfig, c context.Context) (string, error) {
	oconf := makeOAuth2Config(conf)

	if req.FormValue("state") != state {
		return "", errors.WithStack(errors.New("OAuth state doesn't match."))
	}

	token, err := oconf.Exchange(c, req.FormValue("code"))
	if err != nil {
		return "", err
	}
	return stringifyToken(token), nil
}

// GetHTTPClient return an HTTP client that sends authorized requests
func (o OAuth2) GetHTTPClient(token string, conf interactors.ProtocolConfig) (*http.Client, error) {
	oconf := makeOAuth2Config(conf)
	t, err := tokenFromString(token)
	if err != nil {
		return http.DefaultClient, err
	}
	return oconf.Client(context.TODO(), t), nil
}

func makeOAuth2Config(config interactors.ProtocolConfig) *oauth2.Config {
	o := oauth2.Config{
		ClientID:     config.Consumer,
		ClientSecret: config.Secret,
		Endpoint: oauth2.Endpoint{
			TokenURL: config.Endpoint.AccessTokenURL,
			AuthURL:  config.Endpoint.AuthorizationURL,
		},
		RedirectURL: config.CallbackURL,
		Scopes:      config.Scopes,
	}
	return &o
}

func stringifyToken(t *oauth2.Token) string {
	exp, err := t.Expiry.MarshalText()
	if err != nil {
		// The MarshalText function only returns an error if the year of the time is outside the range [0, 9999],
		// which, unless we have a provider acting strangely, shouldn't happen. If we get reports of this
		/// occurring, we'll investigate how to mitigate this.
		panic(err)
	}
	return strings.Join([]string{t.AccessToken, t.RefreshToken, t.TokenType, string(exp)}, "\n")
}

func tokenFromString(s string) (*oauth2.Token, error) {
	parts := strings.Split(s, "\n")

	if len(parts) != 4 {
		return nil, errors.WithStack(errors.New("Incorrect token format, expectted 4 lines oftext."))
	}

	exp := new(time.Time)
	exp.UnmarshalText([]byte(parts[3]))
	t := oauth2.Token{
		AccessToken:  parts[0],
		RefreshToken: parts[1],
		TokenType:    parts[2],
		Expiry:       *exp,
	}
	return &t, nil
}
