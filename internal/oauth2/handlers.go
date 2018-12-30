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

func getProvider(ctx context.Context) (*OAuth, error) {
	pname := ctx.Params().GetString("provider")
	provider := actions.Registry.(*actions2.ProviderRegistry).Provider(pname)

	if provider == nil {
		return nil, errors.New("Provider not found")
	}

	if o, ok := provider.Requirements[requirementName]; ok {
		return o.(*OAuth), nil
	}

	return nil, errors.New("Provider does not use OAuth2.")
}

func CallbackHandler(ctx context.Context) {
	o, err := getProvider(ctx)

	state := ctx.FormValue("state")
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

	conf := o.GoConfig(strings.Split(ctx.Params().GetString("scopes"), ","))
	tok, err := conf.Exchange(ctx.Request().Context(), ctx.FormValue("code"))
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

	err = SaveAccount(o.ProviderName, acc)
	if err != nil {
		ctx.StatusCode(500)
		ctx.WriteString("There was an issue storing the account: " + err.Error())
		return
	}

	ctx.JSON(acc)
}

func Redirect(ctx context.Context) {
	o, err := getProvider(ctx)
	if err != nil {
		ctx.StatusCode(404)
		ctx.WriteString(err.Error())
	}

	state := randx.String(32)
	sess := auth.GetSession(ctx)
	sess.Set("oauth_state", state)

	scopes := strings.Split(ctx.Params().GetString("scopes"), ",")
	conf := o.GoConfig(scopes)
	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
	ctx.Redirect(url)
}
