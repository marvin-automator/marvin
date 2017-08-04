package grifts

import (
	"fmt"
	. "github.com/markbates/grift/grift" // nolint
	"github.com/marvin-automator/marvin/accounts/interactors"
	"github.com/marvin-automator/marvin/accounts/storage"
	globalstorage "github.com/marvin-automator/marvin/storage"
	"strings"
	"github.com/marvin-automator/marvin/accounts/domain"
)

func init() {
	err1 := Desc("accounts:create", "Create a new account.\n $ marvin run accounts:create <email> <password>")
	err2 := Desc("accounts:delete", "Delete an account.\n $ marvin run accounts:delete <id_or_email>")
	if err1 != nil {
		panic(err1)
	}
	if err2 != nil {
		panic(err2)
	}
}

// Grift to create a new account.
var _ = Add("accounts:create", func(c *Context) error {
	globalstorage.Setup()
	defer globalstorage.CloseDB()

	dbs := globalstorage.NewStore()
	defer dbs.Close()
	s := storage.NewAccountStore(dbs)

	action := interactors.CreateAccount{s}
	acc, err := action.Execute(c.Args[0], c.Args[1])
	if err != nil {
		return err
	}

	fmt.Printf("Created account \nEmail: %v\nID: %v\n", acc.Email, acc.ID)
	return nil
})

// Grift to delete an account
var _ = Add("accounts:delete", func(c *Context) error {
	var err error

	globalstorage.Setup()
	defer globalstorage.CloseDB()

	dbs := globalstorage.NewStore()
	defer dbs.Close()
	s := storage.NewAccountStore(dbs)

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

// Grift to list all accounts.
var _ = Add("accounts:list", func(c *Context) error {
	globalstorage.Setup()
	defer globalstorage.CloseDB()

	dbs := globalstorage.NewStore()
	defer dbs.Close()
	s := storage.NewAccountStore(dbs)

	any := false
	fmt.Println("Accounts:")
	s.EachAccount(func(a domain.Account) error {
		any = true
		fmt.Printf("%s: %s", a.Email, a.ID)
		return nil
	})

	if !any {
		fmt.Println("There are no accounts saved. You can create one by running accounts:create.")
	}
	return nil
})