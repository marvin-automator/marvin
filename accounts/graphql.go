package accounts

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/handlers"
)

// GraphQLAccount is the GraphQL representation of an Account.
var GQLAccount = graphql.NewObject(graphql.ObjectConfig{
	Name: "Account",
	Fields: graphql.Fields{
		"id": &graphql.Field{Type: graphql.ID},
		"email": &graphql.Field{Type: graphql.String},
	},
})

var GQLCurrentAccount = graphql.Field{
	Name: "currentAccount",
	Type: GraphQLAccount,
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return CurrentAccount(p.Context.(handlers.Context)), nil
	},
}
