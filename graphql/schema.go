package graphql

import "github.com/graphql-go/graphql"

func schema() *graphql.Schema {
	fields := graphql.Fields{
		"test": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "1, 2, 3", nil
			},
		},
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