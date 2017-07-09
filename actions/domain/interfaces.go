package domain

import (
	"github.com/bigblind/marvin/handlers"
	"github.com/xeipuuv/gojsonschema"
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
type ProviderMeta struct {
	Name        string
	Description string
	// The key should uniquely identify the provider
	Key string
}

// Metadata about a specific action
type ActionMeta struct {
	ProviderMeta
	// The key is used to retrieve the actual action object.
	Key string
	// Whether this action is a trigger
	IsTrigger bool
	// Whether this action needs to do a test run to get the output schema
	RequiresTestRun bool
}

// Returns itself. This is so that action
// implementations can just embed this as a struct
// to get their Meta() method.
func (a ActionMeta) Meta() ActionMeta {
	return a
}

// An action encapsulates how to perform a certain task
type Action interface {
	// Meta returns metadata about the action
	Meta() ActionMeta

	// Called with any data from the front end that is needed to set up the action.
	// The data is passed as raw JSON.
	Setup(data string, c ActionContext) error

	// Should return a struct of the type that this action expects
	// as input.
	// Note that this method is not called on Trigger actions, as
	// they don't get input.
	InputType(c ActionContext) interface{}
	// Actually run the action
	// If RequiresTestRun is true on the Meta, the
	// first call to this method will be a test call.
	// The Context has an IsTestCall method that returns
	// true in this call, and false in every other call.
	Execute(input interface{}, c ActionContext) error
	// The struct type of the data that this action will output.
	OutputType(c ActionContext) interface{}
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
	// The InstanceStore is used to store data that's
	// specific to the current instance of the action.
    // This should be used to store the configuration settings
	// for this particular step in a chore.
	InstanceStore() Store
	// GlobalStore is a Store that provides data that can be reused in any invocation of the action
	// It's the action's responsibility to keep this store clean.
	// The user can clear this store as well.
	GlobalStore() Store
	// Returns a callbackURL that can be used to receive information from other services on the internet
	// The URL is tied to the current instance of the action.
	GetCallbackURL(state string) string
	// Whether this call is a test call.
	// This will only return true in the first call to Execute() for instances of actions
	// where RequiresTestCall is true.
	IsTestCall() bool
	// An action should call Output() to pass output on to the next action.
	// You can call this as many times as you like to provide multiple outputs.
	// In this case, the next step will be called multiple times.
	Output(interface{})
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
