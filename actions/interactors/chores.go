package interactors

import "github.com/bigblind/marvin/actions/domain"

// GetChores is an interactor for retrieving chores
type GetChores struct {
	ChoreStore domain.ChoreStore
}

// ForAccount returns the chores owned by account with ID aid.
func (c GetChores) ForAccount(aid string) ([]domain.Chore, error) {
	cs, err := c.ChoreStore.GetAccountChores(aid)
	if cs == nil {
		cs = []domain.Chore{}
	}
	return cs, err
}