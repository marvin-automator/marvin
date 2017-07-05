package grifts

import (
	"fmt"
	"strings"
	"github.com/bigblind/marvin/actions"
	"github.com/bigblind/marvin/storage"
	. "github.com/markbates/grift/grift"
)

var _ = Add("accounts:create", func(c *Context) error {
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

var _ = Add("accounts:delete", func(c *Context) error {
	storage.Setup()
	defer storage.CloseDB()

	s, err := storage.NewWritableStore()
	defer s.Close()

	if err != nil {
		panic(err)
	}

	action := actions.DeleteAccount{s}
	arg := c.Args[0]
	var t string

	if strings.Contains(arg, "@") {
		err = action.ByEmail(arg)
		t = "email"
	} else {
		err = action.ByID(arg)
		t = "ID"
	}

	if err != nil {
		return err
	}

	fmt.Printf("Deleted account with %v: %v\n", t, arg)
	return nil
})


