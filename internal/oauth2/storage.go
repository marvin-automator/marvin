package oauth2

import (
	"fmt"
	"github.com/marvin-automator/marvin/internal/db"
)

const storeName = "oauth2_accounts"

func makeAccountKey(providerName, accountId string) string{
	return fmt.Sprintf("%v|%v", providerName, accountId)
}

func SaveAccount(providerName string, account Account) error {
	s := db.GetStore(storeName)

	return s.Set(makeAccountKey(providerName, account.Id), account)
}

func GetAccount(providerName, accountId string) (Account, error) {
	s := db.GetStore(storeName)

	a := Account{}
	err := s.Get(makeAccountKey(providerName, accountId), &a)
	return a, err
}