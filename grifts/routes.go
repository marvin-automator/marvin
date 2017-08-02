package grifts

import (
	"fmt"
	"os"
	"text/tabwriter"

	. "github.com/markbates/grift/grift" // nolint
	"github.com/marvin-automator/marvin/app"
)

var _ = Add("routes", func(c *Context) error {
	a := app.App()
	routes := a.Routes()
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
	fmt.Fprintln(w, "METHOD\t PATH\t NAME\t HANDLER")
	fmt.Fprintln(w, "------\t ----\t ----\t -------")
	for _, r := range routes {
		fmt.Fprintf(w, "%s\t %s\t %s\t %s\n", r.Method, r.Path, r.PathName, r.HandlerName)
	}
	err := w.Flush()
	return err
})
