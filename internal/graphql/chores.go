package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/internal/chores"
)

var choreType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Chore",
	Fields: graphql.BindFields(chores.Chore{}),
	Description: "Describes a workflow executed by Marvin",
})

var choreTemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ChoreTemplate",
	Fields: graphql.BindFields(chores.ChoreTemplate{}),
	Description: "Describes a template for chores.",
})

func getChoreQueryFields() graphql.Fields {

	choreType.AddFieldConfig("template", &graphql.Field{
		Type: choreTemplateType,
	})

	return graphql.Fields{
		"Chores": &graphql.Field{
			Type: graphql.NewList(choreType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return chores.GetChores()
			},
		},

		"ChoreById": &graphql.Field{
			Type: choreType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return chores.GetChore(p.Args["id"].(string))
			},
		},

		"ChoreTemplates": &graphql.Field{
			Type: graphql.NewList(choreTemplateType),
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return chores.GetChoreTemplates()
			},
		},

		"ChoreTemplateById": &graphql.Field{
			Type: choreTemplateType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return chores.LoadChoreTemplate(p.Args["id"].(string))
			},
		},
	}
}
