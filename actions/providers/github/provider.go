package github

import (
	"github.com/marvin-automator/marvin/actions/domain"
	"golang.org/x/oauth2/github"
)

func init() {
	p := domain.NewProvider("github", "GitHub", "Actions related to the GitHub code hosting platform")
	p.SetIdentityParameters(github.Endpoint.AuthURL, github.Endpoint.TokenURL, "")

	pt := newPushTrigger()
	p.Add(pt)
}
