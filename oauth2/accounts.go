package oauth2

import "net/http"

type AccountGetter func(client *http.Client) (Account, error)

type Account struct {
	Id       string        `json:"id"`
	Name     string        `json:"name"`
	ImageURL string        `json:"image_url"`
}
