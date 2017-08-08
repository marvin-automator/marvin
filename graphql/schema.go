package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/accounts"
)

func schema() *graphql.Schema {
	fields := graphql.Fields{
		"currentAccount": &accounts.GQLCurrentAccount,
	}

	rq := graphql.ObjectConfig{
		Name: "rootQuery",
		Fields: fields,
	}

	sch, err := graphql.NewSchema(graphql.SchemaConfig{
		Query: rq,
	})

	if err != nil {
		panic(err)
	}

	return &sch
}