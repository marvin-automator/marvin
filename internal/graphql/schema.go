package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var schema *graphql.Schema

func GetHandler() (*handler.Handler, error) {
	if schema == nil {
		_schema, err := getSchema()
		if err != nil {
			return nil, err
		}
		schema = &_schema
	}

	return handler.New(&handler.Config{
		Schema:     schema,
		Playground: true,
		GraphiQL:   true,
		Pretty:     true,
	}), nil
}

func getSchema() (graphql.Schema, error) {
	qt := getQueryType()
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: qt,
	})
}

func getQueryType() *graphql.Object {
	f := combineFields(getChoreQueryFields())
	return graphql.NewObject(graphql.ObjectConfig{
		Name:   "Query",
		Fields: f,
	})
}

func combineFields(f ...graphql.Fields) graphql.Fields {
	res := make(graphql.Fields)

	for _, fields := range f {
		for name, field := range fields {
			res[name] = field
		}
	}

	return res
}
