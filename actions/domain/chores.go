package domain

import (
	"errors"
	"time"
)

var ChoreNotFoundError = errors.New("Chore not found")

// An ActionInstance is an instance of an action in a Chore.
type ActionInstance struct {
	ID             string
	ActionProvider string
	Action         string
	InputTemplate  string
}

// A chore is a workflow specified as a list of actions.
type Chore struct {
	ID      string
	Name    string
	Actions []ActionInstance
	Created time.Time
}

// ChoreStore is an interface for persisting chores.
type ChoreStore interface {
	SaveChore(aid string, c Chore) error
	GetChore(aid, cid string) (Chore, error)
	GetAccountChores(aid string) ([]Chore, error)
	DeleteChore(aid, cid string) error
	DeleteAccountChores(aid string) error
}
