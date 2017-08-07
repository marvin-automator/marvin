package graphql

import (
	"github.com/graphql-go/handler"
	"github.com/bigblind/buffalo"
)

var Handler = buffalo.WrapHandler(handler.New(&handler.Config{
	Schema: schema(),
	Pretty: true,
}))
