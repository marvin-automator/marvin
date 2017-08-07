package accounts

import "github.com/graphql-go/graphql"

// GraphQLAccount is the GraphQL representation of an Account.
var GraphQLAccount = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"email": &graphql.Field{Type: graphql.String},
	},
})

