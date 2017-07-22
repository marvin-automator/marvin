package domain

import "net/http"

type Identity struct {
	Name string
	ImageURL string
	// Subtext can be used for other information that might be helpful to users when selecting an account.
	Subtext string
	// ProviderID is the ID of this identity, assigned by the identity provider
	ProviderID string
	// The access token obtained from the provider
	Token string
}

// IdentityFetcher should be implemented by ActionProviders that want to use
// identity protocols from this package. They're responsible for gathering the account details
// that are not identical across implementations of the protocols
type IdentityFetcher interface{
	// Given an http.Client that will automatically authorize requests,
	// FetchIdentity should fetch the identity associated with the credentials.
	FetchIdentity(c *http.Client) Identity
}

// An IdentityStore persists 3rd party identities
type IdentityStore interface {
	SaveIdentity(account, provider string, id Identity) error
	GetIdentity(account, provider, id string) (Identity, error)
	GetAccountIdentitiesForProvider(account, provider string) []Identity
	DeleteIdentity(account, provider string, i Identity) error
}
