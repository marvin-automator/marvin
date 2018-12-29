package auth

import (
	"bytes"
	"errors"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/sessions"
	"github.com/marvin-automator/marvin/internal/db"
	"html/template"
	"strings"
)

var html, tplErr = packr.NewBox("./templates").FindString("password.html")
var tpl *template.Template

func init() {
	if tplErr != nil {
		panic(tplErr)
	}
	tpl = template.Must(template.New("auth").Parse(html))
}

func renderPasswordTemplate(ctx context.Context, error error, password, next string) {
	w := bytes.NewBufferString("")
	err := tpl.Execute(w, map[string]interface{}{
		"error":    error,
		"password": password,
		"next":     next,
	})

	if err != nil {
		ctx.Writef("Error rendering template: %v", err)
	}
	ctx.HTML(w.String())
}

func AuthHandlers(p iris.Party) {
	p.Use(ensureAuthSession)
	p.Get("/", func(ctx context.Context) {
		next := validateRedirectURL(ctx.FormValueDefault("next", "/"))
		if IsAuthenticated(ctx) {
			ctx.Redirect(next)
			return
		}

		renderPasswordTemplate(ctx, nil, "", next)
	})

	p.Post("/", func(ctx context.Context) {
		next := ctx.FormValueDefault("next", "/")
		pw := ctx.FormValue("password")
		correct, err := IsPasswordValid(pw)

		if err != nil {
			ctx.Writef("An error occurred when checking the password: %v", err)
			return
		}

		if !correct {
			renderPasswordTemplate(ctx, errors.New("Incorrect password"), pw, next)
			return
		}

		s := GetSession(ctx)
		s.Set("authenticated", true)
		ctx.Redirect(next)
	})

	p.Get("/logout", func(ctx context.Context) {
		s := GetSession(ctx)
		s.Set("authenticated", false)
		ctx.Redirect("/", iris.StatusTemporaryRedirect)
	})
}

// Guard against open redirect attacks by making sure the url redirects to the same site, by making sure they start with a slash, and only one slash.
func validateRedirectURL(u string) string {
	if strings.HasPrefix(u, "/") && !strings.HasPrefix(u, "//") {
		return u
	}
	return "/"
}

var sess *sessions.Sessions

func ensureAuthSession(ctx context.Context) {
	if sess == nil {
		if err := makeSessions(); err != nil {
			panic(err)
		}
	}

	s := sess.Start(ctx)
	sess.ShiftExpiration(ctx)

	ctx.Values().Set("session", s)
	ctx.Next()
}

// Returns the session for this request.
func GetSession(ctx context.Context) *sessions.Session {
	return ctx.Values().Get("session").(*sessions.Session)
}

var RequireLogin = context.Handlers{
	ensureAuthSession,
	func(ctx context.Context) {
		if !IsAuthenticated(ctx) {
			next := ctx.Path()
			ctx.Redirect("/auth?next=" + next)
		} else {
			ctx.Next()
		}
	},
}

func IsAuthenticated(ctx context.Context) bool {
	return GetSession(ctx).GetBooleanDefault("authenticated", false)
}

func makeSessions() error {
	hash_key := make([]byte, 64)
	block_key := make([]byte, 32)

	store := db.GetStore("auth")
	if err := store.Get(session_hash_store_key, &hash_key); err != nil {
		return err
	}
	if err := store.Get(session_block_store_key, &block_key); err != nil {
		return err
	}

	s := securecookie.New(hash_key, block_key)

	sess = sessions.New(sessions.Config{
		Cookie:       "auth",
		Encode:       s.Encode,
		Decode:       s.Decode,
		AllowReclaim: true,
	})

	sess.UseDatabase(getSessionDB())

	return nil
}
