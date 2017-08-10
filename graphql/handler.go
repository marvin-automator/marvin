package graphql

import (
	"github.com/graphql-go/handler"
	"github.com/marvin-automator/marvin/handlers"
)

var sch = schema()
var gqlHandler = handler.New(&handler.Config{sch, true})
var marvinHandler handlers.Handler = func(c handlers.Context) error {
	gqlHandler.ContextHandler(c, c.Response(), c.Request())
	return nil
}

// Handler handles GraphQL requests
var Handler = marvinHandler.ToBuffalo()
