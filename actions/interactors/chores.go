package interactors

import "github.com/bigblind/marvin/actions/domain"

type Chores struct {
	ChoreStore domain.ChoreStore
}

func (c Chores) GetForAccount(aid string) ([]domain.Chore, error) {
	cs, err := c.ChoreStore.GetAccountChores(aid)
	if cs == nil {
		cs = []domain.Chore{}
	}
	return cs, err
}