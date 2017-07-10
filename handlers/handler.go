package handlers

import (
	"github.com/bigblind/marvin/storage"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

type Context struct {
	buffalo.Context
}

// Executes the given function with a writable Store instance.
// The store is automatically closed when the function returns.
// If the function returns an error, any changes made to the store are rolled back.
func (c *Context) WithWritableStore(f func(storage.Store) error) error {
	s, err := storage.NewWritableStore()
	if err != nil {
		return err
	}

	err = f(s)
	if err != nil {
		s.RollBack()
	} else {
		s.Close()
	}

	return err
}

func (c Context) WithReadableStore(f func(storage.Store) error) error {
	s, err := storage.NewReadOnlyStore()
	if err != nil {
		return err
	}

	err = f(s)
	s.Close()

	return err
}

func (c Context) Renderer() *render.Engine {
	return r
}

func (c Context) BareRenderer() *render.Engine {
	return br
}

// A request handler that takes our custom context
type Handler func(Context) error

// Turns returns a Buffalo handler that calls this one
func (h Handler) ToBuffalo() buffalo.Handler {
	return func(bc buffalo.Context) error {
		c := Context{bc}
		return h(c)
	}
}
