package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/accounts"
	"github.com/marvin-automator/marvin/actions/actions_graphql"
)

func schema() *graphql.Schema {
	fields := graphql.Fields{
		"currentAccount": &accounts.GQLCurrentAccount,
		"providers": actions_graphql.ProvidersField,
		"viewerChores": actions_graphql.ViewerChoresField,
		"groups": actions_graphql.GroupsField,
	}

	rq := graphql.ObjectConfig{
		Name: "rootQuery",
		Fields: fields,
	}

	sch, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(rq),
	})

	if err != nil {
		panic(err)
	}

	return &sch
}