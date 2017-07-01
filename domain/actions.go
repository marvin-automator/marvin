package domain

import (
	"github.com/xeipuuv/gojsonschema"
)


type ActionProvider interface {
	getActionSet() ActionSet
}

type ActionSet interface {
	getMeta() ActionSetMeta
	getAvailableActions() []ActionMeta
	getAction(key string) Action
}


type ActionSetMeta struct {
	Name string
	Description string
}


type ActionMeta struct {
	ActionSetMeta
	key string
}

type Action interface {
	GetMeta() ActionMeta
	HandleSetup(data string, c ActionContext) error
	GetInputSchema(c ActionContext) gojsonschema.Schema
	GetOutputSchema(c ActionContext) gojsonschema.Schema
	Execute(input string, c ActionContext) error
}

type ActionContext interface{}