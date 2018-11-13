package auth

import (
	"github.com/gorilla/securecookie"
	"github.com/marvin-automator/marvin/internal/db"
	"golang.org/x/crypto/bcrypt"
)

//store keys
const (
	session_hash_store_key  = "session_hash_key"
	session_block_store_key = "session_block_key"
	password_key            = "password"
)

func SetPassword(pw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	s := db.GetStore("auth")
	err = s.Set(password_key, hash)
	if err != nil {
		return err
	}

	// The hash key and block key generated here are used to authenticate and encrypt sessions.
	// Changing them here ensures that any sessions become invalid when the password is changed.
	err = s.Set(session_hash_store_key, securecookie.GenerateRandomKey(64))
	if err != nil {
		return err
	}

	return s.Set(session_block_store_key, securecookie.GenerateRandomKey(32))
}

func IsPasswordValid(pw string) (bool, error) {
	s := db.GetStore("auth")
	var hash []byte
	err := s.Get(password_key, &hash)
	if err != nil {
		return false, err
	}

	if err = bcrypt.CompareHashAndPassword(hash, []byte(pw)); err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	} else {
		return err == nil, err
	}
}

func IsPasswordSet() (bool, error) {
	s := db.GetStore("auth")
	err := s.Get(password_key, &[]byte{})
	if _, ok := err.(db.KeyNotFoundError); ok {
		return false, nil
	}

	return err == nil, err
}
