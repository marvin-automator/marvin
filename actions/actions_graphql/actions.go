package actions_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/actions/domain"
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
		"name": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(domain.Group).Name(), nil
			},
		},
		"actions": &graphql.Field{
			Type: graphql.NewList(Action),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(domain.Group).Actions(), nil
			},
		},
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
				k := p.Source.(domain.ProviderMeta).Key
				return domain.Registry.Provider(k).Groups(), nil
			},
		},
	},
})

var ProvidersField = &graphql.Field{
	Name: "provider",
 	Type: graphql.NewList(Provider),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return domain.Registry.Providers(), nil
	},
}

