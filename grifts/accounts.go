package grifts

import (
	"fmt"
	"github.com/bigblind/marvin/actions"
	"github.com/bigblind/marvin/storage"
	. "github.com/markbates/grift/grift"
)

var _ = Add("create:account", func(c *Context) error {
	storage.Setup()
	defer storage.CloseDB()

	s, err := storage.NewWritableStore()
	defer s.Close()

	if err != nil {
		panic(err)
	}

	action := actions.CreateAccount{s}
	acc, err := action.Execute(c.Args[0], c.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Created account \nEmail: %v\nID: %v\n", acc.Email, acc.ID)
	return nil
})
