package actions_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/actions/interactors"
)


var Action = graphql.NewObject(graphql.ObjectConfig{
	Name: "Action",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"key": &graphql.Field{Type: graphql.ID},
		"isTrigger": &graphql.Field{Type: graphql.Boolean},
	},
})

var ActionGroup = graphql.NewObject(graphql.ObjectConfig{
	Name: "ActionGroup",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
		"actions": &graphql.Field{Type: graphql.NewList(Action)},
		"triggers": &graphql.Field{Type: graphql.NewList(Action)},
	},
})


var Provider = graphql.NewObject(graphql.ObjectConfig{
	Name: "ActionProvider",
	Fields: graphql.Fields{
		"name": &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"key": &graphql.Field{Type: graphql.ID},
		"actionGroups": &graphql.Field{
			Type: graphql.NewList(ActionGroup),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(interactors.Provider).Groups, nil
			},
		},
	},
})

var ProvidersField = &graphql.Field{
	Name: "providers",
 	Type: graphql.NewList(Provider),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return interactors.NewRegistryInteractor().GetProviders(), nil
	},
}

var GroupsField = &graphql.Field{
	Name: "groups",
	Type: graphql.NewList(ActionGroup),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return interactors.NewRegistryInteractor().GetActionGroups(), nil
	},
}

