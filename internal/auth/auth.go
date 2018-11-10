package auth

import (
	"github.com/marvin-automator/marvin/internal/db"
	"golang.org/x/crypto/bcrypt"
)

func SetPassword(pw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s := db.GetStore("auth")
	return s.Set("password", hash)
}

func IsPasswordValid(pw string) (bool, error) {
	s := db.GetStore("auth")
	var hash []byte
	err := s.Get("password", &hash)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword(hash, []byte(pw)); err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return err == nil, err
	}
}
