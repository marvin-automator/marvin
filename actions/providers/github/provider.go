package github

import (
	"github.com/marvin-automator/marvin/actions/domain"
	"golang.org/x/oauth2/github"
	"github.com/gobuffalo/packr"
)

func init() {
	icons := packr.NewBox("./icons")

	p := domain.NewProvider("github", "GitHub", "Actions related to the GitHub code hosting platform")
	p.SetIdentityParameters(github.Endpoint.AuthURL, github.Endpoint.TokenURL, "")
	p.SetIcon(icons.Bytes("github.svg"))

	pt := newPushTrigger()
	pt.SVGIcon = icons.Bytes("commit.svg")
	p.Add(pt)
}
