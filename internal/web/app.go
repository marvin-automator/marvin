package web

import (
	"errors"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"github.com/marvin-automator/marvin/internal/auth"
	"github.com/marvin-automator/marvin/internal/config"
	"github.com/marvin-automator/marvin/internal/graphql"
)

func RunApp() error {
	app := iris.Default()

	set, err := auth.IsPasswordSet()
	if !set {
		return errors.New("set a password using 'marvin set_password' before starting the server")
	}

	if err != nil {
		return err
	}

	app.Any("/", func(ctx context.Context) {
		ctx.Redirect("/app", iris.StatusMovedPermanently)
	})

	app.PartyFunc("/auth", auth.AuthHandlers)
	app.PartyFunc("/app", func(p router.Party) {
		p.Use(auth.RequireLogin...)

		p.Get("/", func(ctx context.Context) {
			ctx.WriteString("Hey, you've logged in!")
		})

		p.Any("/graphql", func(ctx context.Context) {
			h, err := graphql.GetHandler()
			if err != nil {
				ctx.Writef("There was a problem setting up the GraphQL server %v", err)
				return
			}

			h.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		})
	})


	return app.Run(iris.Addr(config.ServerHost))
}
