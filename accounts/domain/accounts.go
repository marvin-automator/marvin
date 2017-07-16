package domain

import (
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrAccountNotFound is returned by functions that can return an account, when the wanted account could not be found.
var ErrAccountNotFound = errors.New("Account not found")

// Account stores information about a user account.
type Account struct {
	ID           string
	Email        string
	PasswordHash []byte
}

// NewAccount creates a new account instance.
func NewAccount(email, password string) (Account, error) {
	id := uuid.NewV4().String()
	hash, err := hashPw(password)
	if err != nil {
		return Account{}, err
	}

	return Account{id, email, hash}, nil
}

// hashPw hashes the password, and returns the hash as bytes.
func hashPw(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

// CheckPassword checks whether the given password equals the user's password, and returns the result as a boolean.
// The error is only set when something went wrong while checking the password.
func (a Account) CheckPassword(pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(pw))

	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}

// The AccountStore interface is an interface for handling the persistence of accounts.
type AccountStore interface {
	SaveAccount(account Account) error
	GetAccountByID(aid string) (Account, error)
	GetAccountByEmail(email string) (Account, error)
	EachAccount(func(a Account) error) error
	DeleteAccount(id string) error
}
