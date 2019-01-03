package oauth2

import "golang.org/x/oauth2"

type AccountGetter func(t *oauth2.Token) (Account, error)

type Account struct {
	Id       string        `json:"id"`
	Token    *oauth2.Token `json:"-"`
	Scopes   []string      `json:"scopes"`
	Name     string        `json:"name"`
	ImageURL string        `json:"image_url"`
}
