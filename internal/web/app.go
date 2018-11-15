package web

import (
	"errors"
	"fmt"
	"github.com/gobuffalo/packr"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/kataras/iris/core/router"
	"github.com/marvin-automator/marvin/internal/auth"
	"github.com/marvin-automator/marvin/internal/config"
	"github.com/marvin-automator/marvin/internal/graphql"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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

	app.PartyFunc("/auth", auth.AuthHandlers)
	app.PartyFunc("/", func(p router.Party) {
		p.Use(auth.RequireLogin...)

		p.Any("/graphql", func(ctx context.Context) {
			h, err := graphql.GetHandler()
			if err != nil {
				ctx.Writef("There was a problem setting up the GraphQL server %v", err)
				return
			}

			h.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
		})

		// In development mode; just pass requests through to the React dev server.
		if config.DevMode {
			cmd := exec.Command("yarn", "start")
			cmd.Dir = "./frontend"
			err := cmd.Start()
			p.HandleMany(iris.MethodGet, "/ /{path:path}", func(ctx context.Context) {
				if err != nil {
					ctx.Writef("Error starting dev server %v.", err)
					return
				}

				url := "http://localhost:3000/" + ctx.Params().GetString("path")
				fmt.Printf("Loading from %v\n", url)
				resp, err := http.Get(url)
				if err != nil {
					ctx.Writef("Error Getting page: %v.", err)
					return
				}

				ctx.StatusCode(resp.StatusCode)
				for h, v := range resp.Header {
					ctx.Header(h, strings.Join(v, ","))
				}

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					ctx.Writef("Error reading body %v", err)
				}

				ctx.Write(body)
			})
		} else {
			// Otherwise, we use packr to bundle in the built frontend
			box := packr.NewBox("../../frontend/build")

			p.Get("/{path:path}", func(ctx context.Context) {
				f, err := box.Open(ctx.Params().GetString("path"))
				if err == os.ErrNotExist {
					f, err = box.Open("index.html")
				}

				if err != nil {
					ctx.StatusCode(500)
					ctx.Writef("Error getting file: %v", err)
					return
				}

				info, err := f.Stat()

				if err != nil {
					ctx.StatusCode(500)
					ctx.Writef("Error getting file: %v", err)
					return
				}

				ftype := mime.TypeByExtension(filepath.Ext(info.Name()))
				ctx.Header("Content-Type", ftype)
				ctx.Header("Content-Length", strconv.Itoa(int(info.Size())))

				body := make([]byte, info.Size())
				f.Read(body)
				ctx.Write(body)
			})
		}
	})

	return app.Run(iris.Addr(config.ServerHost))
}
