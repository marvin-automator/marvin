package actions

// Requirements specify the things an action needs to be able to run.
type Requirement interface {
	Name() string
	Init(providerName string)
	Config() interface{}
	ConfigHelp() string
	SetConfig(interface{}) error
	Fulfilled() bool
}








