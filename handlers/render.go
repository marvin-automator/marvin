package handlers

import (
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/plush"
	"github.com/tonnerre/golang-pretty"
)

var r *render.Engine
// bare renderer
var br *render.Engine

func DebugHelper(x string, help plush.HelperContext) string {
	return pretty.Sprintf("%# v", help)
}

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),

		// Add template helpers here:
		Helpers: render.Helpers{
			"form": plush.FormHelper,
		},
	})

	br = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "bare.html",

		// Box containing all of the templates:
		TemplatesBox: packr.NewBox("../templates"),

		// Add template helpers here:
		Helpers: render.Helpers{
			"form": plush.FormHelper,
			"showContext": DebugHelper,
		},
	})
}
