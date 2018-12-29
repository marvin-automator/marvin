package oauth2

import (
	"crypto/subtle"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/errors"
	"github.com/markbates/going/randx"
	"github.com/marvin-automator/marvin/actions"
	actions2 "github.com/marvin-automator/marvin/internal/actions"
	"github.com/marvin-automator/marvin/internal/auth"
	"golang.org/x/oauth2"
	"strings"
)

func getProvider(ctx context.Context) (*actions2.Provider, error) {
	pname := ctx.Values().GetString("provider")
	provider := actions.Registry.(*actions2.ProviderRegistry).Provider(pname)

	if provider == nil {
		return nil, errors.New("Provider not found")
	}

	if provider.Requirements.OAuth2 == nil {
		return nil, errors.New("Provider does not use OAuth2.")
	}

	return provider, nil
}

func CallbackHandler(ctx context.Context) {
	p, err := getProvider(ctx)
	o := p.Requirements.OAuth2

	state := ctx.Values().GetString("state")
	sess := auth.GetSession(ctx)
	if subtle.ConstantTimeCompare([]byte(state), []byte(sess.GetString("oauth_state"))) != 1 {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("Incorrect state parameter. This can happen because you ended up on this page accidentally, or your cookies have been deleted/modified.")
	}

	if err != nil {
		ctx.StatusCode(404)
		ctx.WriteString(err.Error())
		return
	}

	conf := o.GoConfig(strings.Split(ctx.Values().GetString("scopes"), ","))
	tok, err := conf.Exchange(ctx.Request().Context(), ctx.Values().GetString("code"))
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString(err.Error())
		return
	}

	acc, err := o.GatAccount(tok)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString("There was an issue getting the account: " + err.Error())
		return
	}

	err = SaveAccount(p.Name, acc)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString("There was an issue storing the account: " + err.Error())
		return
	}

	ctx.JSON(acc)
}

func Redirect(ctx context.Context) {
	p, err := getProvider(ctx)
	if err != nil {
		ctx.StatusCode(404)
		ctx.WriteString(err.Error())
	}

	state := randx.String(32)
	sess := auth.GetSession(ctx)
	sess.Set("oauth_state", state)

	o := p.Requirements.OAuth2
	scopes := strings.Split(ctx.Values().GetString("scopes"), ",")
	ctx.Redirect(o.GoConfig(scopes).AuthCodeURL(state, oauth2.AccessTypeOffline))
}
