package auth

import (
	"github.com/marvin-automator/marvin/internal/db"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSetAndCheckPassword(t *testing.T) {
	r := require.New(t)
	db.SetupTestDB()

	err := SetPassword("badPassword")
	r.NoError(err)

	ok, err := IsPasswordValid("badPassword")

	r.NoError(err)
	r.True(ok, "The password was supplied correctly, so this should be true.")

	db.TearDownTestDB()
}

func TestIncorrectPassword(t *testing.T) {
	r := require.New(t)
	db.SetupTestDB()

	err := SetPassword("badPassword")
	r.NoError(err)

	ok, err := IsPasswordValid("completelyD!fferent")

	r.NoError(err)
	r.False(ok, "The password was supplied incorrectly, so this should be false.")
	db.TearDownTestDB()
}
