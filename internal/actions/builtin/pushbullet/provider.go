package pushbullet

import (
	"encoding/json"
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal/actions/builtin/icons"
	"github.com/marvin-automator/marvin/oauth2"
	oauth22 "golang.org/x/oauth2"
	"net/http"
)

func init() {
	p := actions.Registry.AddProvider("pushbullet", "Send push notifications using PushBullet", icons.Get("pusher.svg"))

	p.AddRequirement(oauth2.Requirement(oauth22.Endpoint{
		AuthURL:  "https://www.pushbullet.com/authorize",
		TokenURL: "https://api.pushbullet.com/oauth2/token",
	}, helpTemplate, getAccount))
}

func getAccount(client *http.Client) (oauth2.Account, error) {
	resp, err := client.Get("https://api.pushbullet.com/v2/users/me")
	if err != nil {
		return oauth2.Account{}, err
	}

	dec := json.NewDecoder(resp.Body)
	acc := account{}
	err = dec.Decode(&acc)

	return oauth2.Account{
		Name:     acc.Name,
		Id:       acc.Iden,
		ImageURL: acc.ImageURL,
	}, err
}

type account struct {
	Iden     string `json:"iden"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

var helpTemplate = `
<p>
	<a href="https://www.pushbullet.com/#settings/clients">Click here to sign up or log in to PushBullet</a>.
	If this link didn't take you to a page named OAuth Clients, go to "settings" and then click "Clients".
</p>

<p>You can enter whatever name, website_url and image_url you want. For the redirect_uri enter "{{.RedirectURL}}",
and leave the allowed_origin blank. Now click "Add A new OAuth Client".</p>

<p>Enter the client_id and client_secret in the form below.</p>
`
