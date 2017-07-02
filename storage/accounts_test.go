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
