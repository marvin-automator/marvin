package interactors

import (
	"net/http"
	"context"
)

type IdentityProtocol int

const (
	None IdentityProtocol = 0
	OAuth1 = iota
	OAuth2 = iota
)

var implementations = map[IdentityProtocol]Protocol{}

// Register an implementation for the given protocol
func RegisterImplementation(p IdentityProtocol, impl Protocol) {
	implementations[p] = impl
}

func GetProtocol(p IdentityProtocol) Protocol {
	return implementations[p]
}

// Endpoint represents a set of URLs offered by an identity
// provider to handle the autorization protocol.
type Endpoint struct {
	// The URL users are redirected to for authorization
	AuthorizationURL string
	// UrL where an access token can be obtained
	AccessTokenURL string
	// RequestTokenURL is a URL where a temporary token used in oauth 1.0 can be obtained.
	RequestTokenURL string
}

// ProtocolConfig holds configuration values used by protocol implementations
type ProtocolConfig struct {
	// Consumer somehow identifies this application
	// to the authentication provider. Some providers
	// call this the App ID or something similar.
	Consumer string
	// Secret authenticates the application to the auth provider.
	Secret string
	// Endpoint stores endpoints offered by the oauth provider for
	// /handling authentication and authorization.
	Endpoint Endpoint
	// CallbackURL holds the URL that the client should be redirected to after
	// the user has granted authorization.
	CallbackURL string
	// Scopes holdds the scopes to request.
	Scopes []string
}

// Protocol is a protocol for requesting 3rd party identities and authorization
type Protocol interface {
	// AuthorizationURL returns the URL that users should be redirected to to obtain authorization.
	// Scopes will be a comma-separated list of scopes
	AuthorizationURL(conf ProtocolConfig, state string, c context.Context) (string, error)
	// GetToken obtains a token from the authentication provider, given the request sent to the callback URL
	GetToken(req http.Request, state string, conf ProtocolConfig, c context.Context) (string, error)
	// GetHTTPClient returns a HTTP client where requests made using it carry authentication information for the provider
	GetHTTPClient(token string, conf ProtocolConfig) (*http.Client, error)
}
