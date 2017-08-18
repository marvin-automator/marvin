package interactors

import "github.com/marvin-automator/marvin/actions/domain"

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

type TemporaryInstances struct {
	TemporaryInstanceStore domain.TemporaryInstanceStore
}

func (ti TemporaryInstances) New(provider, action string) (domain.TemporaryActionInstance, error) {
	tai := domain.NewTemporaryActionInstance(provider, action)
	return tai, ti.TemporaryInstanceStore.Save(tai)
}

func (ti TemporaryInstances) ByID(ID string) (domain.TemporaryActionInstance, error) {
	return ti.TemporaryInstanceStore.Get(ID)
}