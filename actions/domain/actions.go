package domain

import (
	"github.com/xeipuuv/gojsonschema"
	"context"
	"github.com/bigblind/marvin/handlers"
)



// ActionProvider provides a list of related actions.
type ActionProvider interface {
	// Meta returns metadata about this provider
	Meta() ActionSetMeta

	// ActionList provides a list of the available actions.
	ActionList() []ActionMeta

	// The action method gives access to the interface for actually configuring and running actions.
	Action(key string) Action
}

// Metadata about a set of actions
type ActionSetMeta struct {
	Name        string
	Description string
}

// Metadata about a specific action
type ActionMeta struct {
	ActionSetMeta
	// The key is used to retrieve the actual action object.
	Key string
	// Returns whether this action is a trigger
	IsTrigger bool
}

// An action encapsulates how to perform a certain task
type Action interface {
	// Meta returns metadata about the action
	Meta() ActionMeta

	// Called with any data from the front end that is needed to set up the action.
	// The data is passed as raw JSON.
	Setup(data string, c ActionContext) error

	// Represents the structure of the data that this action expects.
	// This is not the initial configuration input, but the input
	// it gets from previous actions.
	// Note that this method is not called on Trigger actions, as
	// they don't get input.
	InputSchema(c ActionContext) gojsonschema.Schema
	// Actually run the action
	Execute(input string, c ActionContext) error
	// The schema of the data that this action will output.
	OutputSchema(c ActionContext) gojsonschema.Schema
}

// A trigger is an action that starts
type Trigger interface {
	Action
	// Called when marvin is started, for triggers that
	// need to run continuously.
	Start(c ActionContext) error
	// Callback gets called when a callback URL is invoked
	// It receives the state that was passed to GetCallbackURL,
	// and should return a handler to handle the request.
	Callback(state string) handlers.Handler
}

// Gives you access to data and functionality that can be useful when executing an action
type ActionContext interface{
	// InvocationStore is a Store that stores data for the duration of the invocation of this action.
	// The data is automatically deleted when you call Done
	InvocationStore() Store
	// GlobalStore is a Store that provides data that can be reused in any invocation of the action
	// It's the action's responsibility to keep this store clean.
	// The user can clear this store as well.
	GlobalStore() Store
	// Returns a callbackURL that can be used to receive information from other services on the internet
	GetCallbackURL(state string) string
}

// Store is an interface for storing data.
type Store interface {
	// Get data from the store. Returns nil if the key isn't there.
	Get(key string) (interface{}, error)
	// Put the value into the store
	Put(key string, value interface{}) error
	// Delete the ivalue associated with this key from the store.
	Delete(key string) error
}
