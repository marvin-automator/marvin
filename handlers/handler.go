package handlers

import (
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/marvin-automator/marvin/storage"
)

// Context is a marvin-specific context object.
type Context struct {
	buffalo.Context
	store storage.Store
}

// Store returns a Store. The Store should only be used while this context is valid.
func (c Context) Store() storage.Store {
	return c.store
}

// Renderer returns a buffalo rendering engine, configured to use the main application layout file
func (c Context) Renderer() *render.Engine {
	return r
}

// BareRenderer returns a buffallo rendering engine, configured to use the bare html layout file.
func (c Context) BareRenderer() *render.Engine {
	return br
}

// Handler is a request handler that takes our custom context
type Handler func(Context) error

// ToBuffalo turns returns a Buffalo handler that calls this one
func (h Handler) ToBuffalo() buffalo.Handler {
	s := storage.NewStore()

	return func(bc buffalo.Context) error {
		c := Context{bc, s}
		err := h(c)

		go func() {
			<-bc.Done() // When the context is done with; close the store.
			bc.Logger().Debug("Closing store")
			s.Close()
			bc.Logger().Debug("Closed store")
		}()

		return err
	}
}
