package grifts

import (
	"fmt"
	"github.com/bigblind/marvin/accounts/interactors"
	"github.com/bigblind/marvin/accounts/storage"
	globalstorage "github.com/bigblind/marvin/storage"
	. "github.com/markbates/grift/grift"
	"strings"
)

var _ = Add("accounts:create", func(c *Context) error {
	globalstorage.Setup()
	defer globalstorage.CloseDB()

	dbs, err := globalstorage.NewWritableStore()
	defer dbs.Close()
	s := storage.NewAccountStore(dbs)

	if err != nil {
		panic(err)
	}

	action := interactors.CreateAccount{s}
	acc, err := action.Execute(c.Args[0], c.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Created account \nEmail: %v\nID: %v\n", acc.Email, acc.ID)
	return nil
})

var _ = Add("accounts:delete", func(c *Context) error {
	globalstorage.Setup()
	defer globalstorage.CloseDB()

	dbs, err := globalstorage.NewWritableStore()
	defer dbs.Close()
	s := storage.NewAccountStore(dbs)

	if err != nil {
		panic(err)
	}

	action := interactors.DeleteAccount{s}
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
