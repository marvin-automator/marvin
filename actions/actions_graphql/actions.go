package actions_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/actions/domain"
	"github.com/marvin-automator/marvin/actions/interactors"
	"github.com/pkg/errors"
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
				switch t := p.Source.(type) {
				case domain.Group:
					return t.Name(), nil
				case interactors.Group:
					return t.Name, nil
				default:
					return "", errors.New("Unexpected podcast type")
				}
			},
		},
		"actions": &graphql.Field{
			Type: graphql.NewList(Action),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				switch t := p.Source.(type) {
				case domain.Group:
					return t.Actions(), nil
				case interactors.Group:
					return t.Actions, nil
				default:
					return nil, errors.New("Unexpected podcast type")
				}
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
	Name: "providers",
 	Type: graphql.NewList(Provider),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return domain.Registry.Providers(), nil
	},
}

var GroupsField = &graphql.Field{
	Name: "groups",
	Type: graphql.NewList(ActionGroup),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return interactors.NewRegistryInteractor().GetActionGroups(), nil
	},
}

