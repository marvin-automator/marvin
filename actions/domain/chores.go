package domain

import (
	"errors"
	"time"
)

// ErrChoreNotFound is returned by methods that should return a chore when the requested chore cannot be found.
var ErrChoreNotFound = errors.New("chore not found")

// An ActionInstance is an instance of an action in a Chore.
type ActionInstance struct {
	ID             string `json:"id"`
	ActionProvider string `json:"action_provider"`
	Action         string `json:"action"`
	InputTemplate  string `json:"input_template"`
	Identity	   string `json:"identity"`
}

// A Chore is a workflow specified as a list of actions.
type Chore struct {
	ID      string			 `json:"id"`
	Name    string           `json:"name"`
	Actions []ActionInstance `json:"actions"`
	Created time.Time        `json:""`
}

// ChoreStore is an interface for persisting chores.
type ChoreStore interface {
	SaveChore(aid string, c Chore) error
	GetChore(aid, cid string) (Chore, error)
	GetAccountChores(aid string) ([]Chore, error)
	DeleteChore(aid, cid string) error
	DeleteAccountChores(aid string) error
}
