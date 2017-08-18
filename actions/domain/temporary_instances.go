package domain

import "time"

type TemporaryActionInstance struct{
	ActionInstance
	Created time.Time
}

func NewTemporaryActionInstance(provider, action string) *TemporaryActionInstance {
	return &TemporaryActionInstance{
		ActionInstance: NewActionInstance(provider, action),
		Created: time.Now(),
	}
}

type TemporaryActionStore interface{
	Save(tai TemporaryActionInstance) error
	Get(ID string) (TemporaryActionInstance, error)
	DeleteInstancesOlderThan(t time.Time) (int, error)
}