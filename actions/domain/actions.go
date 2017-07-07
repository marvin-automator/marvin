package domain

import (
	"github.com/xeipuuv/gojsonschema"
)

type ActionProvider interface {
	getActionSet() ActionSet
}

type ActionSet interface {
	Meta() ActionSetMeta
	ActionList() []ActionMeta
	Action(key string) Action
}

type ActionSetMeta struct {
	Name        string
	Description string
}

type ActionMeta struct {
	ActionSetMeta
	key string
}

type Action interface {
	Meta() ActionMeta
	Setup(data string, c ActionContext) error
	InputSchema(c ActionContext) gojsonschema.Schema
	OutputSchema(c ActionContext) gojsonschema.Schema
	Execute(input string, c ActionContext) error
}

type ActionContext interface{}
