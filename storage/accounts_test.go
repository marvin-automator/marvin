package storage

import (
	"testing"
	"github.com/bigblind/marvin/domain"
	"bytes"
)

func TestSaveAndGetAccount(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		if err != nil {
			t.Error(err)
		}

		a2, err := s.GetAccountByID("042")
		if err != nil {
			t.Error(err)
		}

		if a1.ID != a2.ID || a1.Email != a2.Email || !bytes.Equal(a1.PasswordHash, a2.PasswordHash){
			t.Error("The saved account does not equal the retrieved account.")
		}

	})
}

func TestGetAccountByEmailExists(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		if err != nil {
			t.Error(err)
		}

		a2, err := s.GetAccountByEmail("foo@example.com")
		if err != nil {
			t.Error(err)
		}

		if a1.ID != a2.ID || a1.Email != a2.Email || !bytes.Equal(a1.PasswordHash, a2.PasswordHash){
			t.Errorf("The saved account does not equal the retrieved account: \n%v\n%v", a1, a2)
		}

	})
}

func TestGetAccountByEmailDoesNotExist(t *testing.T) {
	WithTestDB(t, func(s Store) {
		a1 := domain.Account{"042", "foo@example.com", []byte("nothashed")}
		err := s.SaveAccount(a1)
		if err != nil {
			t.Error(err)
		}

		a2, err := s.GetAccountByEmail("bar@example.com")
		if err != AccountNotFound {
			t.Errorf("Account should not be found, but didn't get AccountNotFound error. got %v instead.", err)
		}
		if a2.ID != "" || a2.Email != "" || !bytes.Equal(a2.PasswordHash, []byte{}) {
			t.Error("The passed in account should not be altered.")
		}
	})
}


