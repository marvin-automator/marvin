package domain

import "time"

// TemporaryActionInstance holds an ActionInstance that's being configured to be put into a Chore
type TemporaryActionInstance struct{
	ActionInstance
	Created time.Time
}

// NewTemporaryActionInstance creates a new TemporaryActionInstance with an action instance for the given provider and action.
func NewTemporaryActionInstance(provider, action string) TemporaryActionInstance {
	return TemporaryActionInstance{
		ActionInstance: NewActionInstance(provider, action),
		Created: time.Now(),
	}
}

//TemporaryInstanceStore stores and retrieves TemporaryActionInstances
type TemporaryInstanceStore interface{
	Save(tai TemporaryActionInstance) error
	Get(ID string) (TemporaryActionInstance, error)
	DeleteInstancesOlderThan(t time.Time) (int, error)
}