package domain

import (
	"errors"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var ErrAccountNotFound = errors.New("Account not found")

type Account struct {
	ID           string
	Email        string
	PasswordHash []byte
}

func NewAccount(email, password string) (Account, error) {
	id := uuid.NewV4().String()
	hash, err := hashPw(password)
	if err != nil {
		return Account{}, err
	}

	return Account{id, email, hash}, nil
}

func hashPw(pwd string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
}

func (a Account) checkPassword(pw string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(a.PasswordHash, []byte(pw))

	if err == nil {
		return true, nil
	} else if err == bcrypt.ErrMismatchedHashAndPassword {
		return false, nil
	}

	return false, err
}

type AccountStore interface {
	SaveAccount(account Account) error
	GetAccountByID(aid string) (Account, error)
	GetAccountbyEmail(email string) (Account, error)
}
