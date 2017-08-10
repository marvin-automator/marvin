package actions_graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/actions/interactors"
	"github.com/marvin-automator/marvin/handlers"
	"github.com/marvin-automator/marvin/actions/storage"
	accountsinteractors "github.com/marvin-automator/marvin/accounts/interactors"

	"github.com/marvin-automator/marvin/actions/domain"
)

var ActionInstance = graphql.NewObject(graphql.ObjectConfig{
	Name: "ActionInstance",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"inputTemplate": &graphql.Field{Type: graphql.String},

		"action": &graphql.Field{
			Type: Action,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				act := p.Source.(domain.ActionInstance)
				return domain.Registry.Provider(act.ActionProvider).Action(act.Action).Meta(), nil
			},
		},
		"provider": &graphql.Field{
			Type: Provider,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				act := p.Source.(domain.ActionInstance)
				return domain.Registry.Provider(act.ActionProvider).Meta(), nil
			},
		},
	},
})


var Chore = graphql.NewObject(graphql.ObjectConfig{
	Name: "Chore",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"name": &graphql.Field{Type: graphql.String},
		"owner": &graphql.Field{Type: graphql.String},
		"actions": &graphql.Field{
			Type: graphql.NewList(Action),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return p.Source.(domain.Chore).Actions, nil
			},
		},
	},
})


var ViewerChoresField = &graphql.Field{
	Name: "viewerChores",
	Type: graphql.NewList(Chore),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		ctx := p.Context.(handlers.Context)
		ctx.Logger().Info("Graphql /viewerChores")
		cs := storage.NewChoreStore(ctx.Store())
		i := interactors.GetChores{cs}
		ctx.Logger().Debug("Calling the interactor")
		chs, err := i.ForAccount(ctx.Value("account").(accountsinteractors.Account).ID)
		ctx.Logger().Debug("Call done")
		return chs, err
	},

}
