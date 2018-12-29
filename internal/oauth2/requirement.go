package oauth2

import (
	"bytes"
	"github.com/kataras/iris/core/errors"
	"github.com/marvin-automator/marvin/actions"
	"golang.org/x/oauth2"
	"text/template"
)

type OAuth struct {
	Provider actions.Provider
	Endpoint oauth2.Endpoint
	ConfigHelpTemplate string
	config Config
	GatAccount AccountGetter
}

func (o *OAuth) Init(p actions.Provider) {
	o.Provider = p
}

func (o *OAuth) ConfigHelp() string {
	wr := bytes.NewBufferString("")
	template.Must(template.New("help_text").Parse(o.ConfigHelpTemplate)).Execute(wr, o)
	return wr.String()
}

func (o *OAuth) Config() interface{} {
	return o.config
}

func (o *OAuth) SetConfig(c interface{}) error {
	conf, ok := c.(Config)
	if !ok {
		return errors.New("OAuth2 config has incorrect type.")
	}

	o.config = conf
	return nil
}

func (o *OAuth) Fulfilled() bool {
	return o.config.ClientID != "" && o.config.Secret != ""
}

func (o *OAuth) GoConfig(scopes []string) oauth2.Config {
	return oauth2.Config{
		Endpoint: o.Endpoint,
		ClientID: o.config.ClientID,
		ClientSecret: o.config.Secret,
		Scopes: scopes,
		RedirectURL: o.RedirectURL(),
	}
}

func (o *OAuth) RedirectURL() string {
	return  "/oauth/callback/" + o.Provider.Info().Name
}

type Config struct {
	ClientID	string 	`json:"client_id"`
	Secret		string 	`json:"secret"`
}

