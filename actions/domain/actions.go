package domain

import (
	"github.com/marvin-automator/marvin/app/domain"
	"net/http"
)

// ActionProvider provides a list of related actions.
type ActionProvider interface {
	// Meta returns metadata about this provider
	Meta() ProviderMeta

	// Groups returns a list of action groups.
	// See the Group type for more info.
	Groups() []Group

	// The action method gives access to the interface for actually configuring and running actions.
	Action(key string) BaseAction

	// GlobalConfigType should return an instance of a struct type, that'll
	// hold global configuration data. This should be used for configuration data
	// that's not account-specific, like API client secrets.
	// If no global configuration is needed, return nil.
	GlobalConfigType() interface{}
}

// Group groups related actions within a provider together. Most action providers will return a
// single group with the same name as the provider,
// but groups provide a way to subcaterogize them. Groups show up as separate providers in the
// action selection list, but share the global data store, so they can share accounts, etc.
type Group interface {
	// Actions returns ActionMetas of actions that are available in this group
	Actions() []ActionMeta

	// Name returns a human-readable name for the group
	Name() string
}

// IdentityProtocol identifies an authentication protocol.
type IdentityProtocol uint8

const (
	None   IdentityProtocol = iota // No identity protocol needed
	OAuth1 IdentityProtocol = iota // OAuth 1.0a
	OAuth2 IdentityProtocol = iota // Oauth 2.0
)

// ProviderMeta stores metadata about a set of actions
type ProviderMeta struct {
	Name        string	`json:"name"`
	Description string  `json:"description"`

	// The key should uniquely identify the provider
	Key 		string	`json:"key"`

	// The protocol that should be used for getting identities to use with the actions for this provider.
	// Use None if you don't need a 3rd-party identity.
	//
	// If your provider does require 3rd-party identities, your GlobalConfigType must have the following fields:
	// ClientID string // Some providers call this a consumer ID, App ID, or other name.
	// ClientSecret string // Some  providers may call this the Consumer Secret or App Secret
	ReequiresIdentityProtocol IdentityProtocol

	// If 3rd-party identities are required, AuthorizationEndpoint is the authorization endpoint of the provider
	AuthorizationEndpoint string

	// If 3rd-party identities are required, TokenEndpoint is the URL where a token can be obtained
	TokenEndpoint string

	// If 3rd-party identities are required, and uses Oauth 1.0, RequestTokenEndpoint is the URL where a request token can be obtained.
	RequestTokenEndpoint string

	// The scopes to require for a 3rd-party identity
	Scopes []string
}

// ActionMeta stores metadata about a specific action
type ActionMeta struct {
	// The key should uniquely identify the action within the provider
	Key         string	`json:"key"`
	Name        string  `json:"name"`
	Description string  `json:"description"`

	// SetupStepPath, should either be an empty string if your action doesn't need a setup step, or a path that your
	// action can handle via the callback URL. The response from this path will be included in an iframe during the
	// action setup process. To indicate that the setup step is complete, call the following function in JavaScript:
	//
	//   window.top.postMessage("*", "done");
	//
	SetupStepPath string

	// The key is used to retrieve the actual action object.
	// Whether this action is a trigger
	IsTrigger bool    `json:"isTrigger"`

	// Whether this action needs to do a test run to get the output schema
	RequiresTestRun bool

	// Whether this action requires a 3rd-party identity
	RequiresIdentity bool
}

// Meta returns itself. This is so that action
// implementations can just embed this as a struct
// to get their Meta() method.
func (a ActionMeta) Meta() ActionMeta {
	return a
}

// SetMeta here is used as follows: Embed an ActionMeta in your action. Then, you can call this method on a pointer to the action to easily set the meta information.
func (a *ActionMeta) SetMeta(key, name, description string, isTrigger, requiresTestRun bool, requiresIdentity bool) {
	a.Key = key
	a.Name = name
	a.Description = description
	a.IsTrigger = isTrigger
	a.RequiresTestRun = requiresTestRun
	a.RequiresIdentity = requiresIdentity
}

// BaseAction is the interface that Trigger and Action have in common
type BaseAction interface {
	// Meta returns metadata about the action
	// To avoid having to implement this, embed an instance of ActionMeta
	Meta() ActionMeta

	// The struct type of the data that this action will output.
	// If RequiresTestCall is true for this action's meta, OutputType will not be called,
	// The action will be run instead and whatever it outputs will be used as its output type.
	OutputType(c ActionContext) interface{}
}

// A Trigger is an action that starts off a chore.
type Trigger interface {
	BaseAction

	// Start starts off the trigger
	Start(c ActionContext)
}

// An Action encapsulates how to perform a certain task
type Action interface {
	BaseAction

	// Actually run the action
	// If RequiresTestRun is true on the Meta,
	// Execute will also be executed for test calls.
	// The Context has an IsTestCall method that returns
	// true in these calls, and false in normal calls.
	Execute(input interface{}, c ActionContext) error

	// Should return a struct of the type that this action expects
	// as input.
	// Note that this method is not called on Trigger actions, as
	// they don't get input.
	InputType(c ActionContext) interface{}
}

// CallbackReceiver should be implemented by actions that want to receive requests from callback URLs.
type CallbackReceiver interface {
	// Callback gets called when a callback URL is invoked. It receives the path, relative to its callback URL.
	Callback(req *http.Request, rw ActionResponseWriter, path string, c ActionContext)
}

// ActionResponseWriter is a http.ResponseWriter-derived interface
// that actions use to respond to callback URLs
type ActionResponseWriter interface {
	http.ResponseWriter

	// Send a json response with the given code. i will be marshalled into JSON.
	JSON(code int, i interface{})

	// Send an HTML response
	HTML(code int, content string)

	// Send a plaintext response
	Text(code int, content string)
}

// ActionContext gives an action access to data and functionality that can be useful when executing an action
type ActionContext interface {
	// InvocationStore is a KVStore that stores data for the duration of the invocation of this action.
	// The data is automatically deleted when you call Done
	InvocationStore() KVStore

	// The InstanceStore is used to store data that's
	// specific to the current instance of the action.
	// This should be used to store the configuration settings
	// for this particular step in a chore. It can also be used to accumulate/change data across invocations.
	InstanceStore() KVStore

	// AccountGlobalStore is a KVStore that provides data that can be reused in any invocation of the action
	// within the current user's account. It's the action's responsibility to keep this store clean.
	// The user can clear this store as well.
	AccountGlobalStore() KVStore

	// HTTPClient returns a http.Client that'll make authenticated requests with the identity, selected for this action.
	// If the action doesn't require a 3rd-party identity, this just returns an unmodified http.Client.
	HTTPClient() *http.Client

	// Returns a callbackURL that can be used to receive information from other services on the internet
	// The URL is tied to the current invocation of the action. The domain and start of the path are fixed, buy
	// you can provide your own path suffix to route things inside your action. The path suffix must always start with a /,
	// or be the empty string.
	GetCallbackURL(path string) string

	// Whether this call is a test call.
	IsTestCall() bool

	// GlobalConfig returns the global config object. It should be convertible to the type
	// returned from Action.GlobalConfigType(), unless that returned nil, in which case this
	// function returns nil.
	GlobalConfig() interface{}

	// An action should call Output() to pass output on to the next action.
	// You can call this as many times as you like to provide multiple outputs.
	// In this case, the next step will be called multiple times.
	Output(interface{})

	// An action should call Done() when it is done sending outputs.
	// A trigger should never call Done, as it should keep sending outputs until the chore is deleted
	Done()

	// Logger returns a Logger instance actions can use to log messages that help users understand what's going on.
	Logger() domain.Logger
}

// KVStore is an interface for storing data in actions.
type KVStore interface {

	// Get data from the store. Returns nil if the key isn't there.
	Get(key string) (interface{}, error)

	// Put the value into the store
	Put(key string, value interface{}) error

	// Delete the ivalue associated with this key from the store.
	Delete(key string) error

	// Close needs to be called to free up resources on the database.
	// Stores are automatically closed when the action that requested them finishes,
	// but in triggers, it's important to close the store when as soon as possible
	Close()
}
