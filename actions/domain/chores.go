package domain

// An ActionInstance is an instance of an action in a Chore.
type ActionInstance struct {
	ID string
	ActionProvider string
	Action string
	InputTemplate string
}

// A chore is a workflow specified as a list of actions.
type Chore struct {
	Name string
	Actions []ActionInstance
	Created time.Time
}
