package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/internal/chores"
)

var idArgs = graphql.FieldConfigArgument{
	"id": &graphql.ArgumentConfig{
		Type: graphql.String,
	},
}

var choreType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "Chore",
	Fields:      graphql.BindFields(chores.Chore{}),
	Description: "Describes a workflow executed by Marvin",
})

var choreTemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ChoreTemplate",
	Fields:      graphql.BindFields(chores.ChoreTemplate{}),
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
			Args: idArgs,
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
			Args: idArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return chores.LoadChoreTemplate(p.Args["id"].(string))
			},
		},
	}
}

func getChoreMutationFields() graphql.Fields {
	return graphql.Fields{
		"createChoreTemplate": &graphql.Field{
			Type: choreTemplateType,
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"script": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := chores.NewChoreTemplate(p.Args["name"].(string), p.Args["script"].(string))
				if err != nil {
					return nil, err
				}

				return c, c.Save()
			},
		},

		"updateChoreTemplate": &graphql.Field{
			Type: choreTemplateType,
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"script": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := chores.LoadChoreTemplate(p.Args["id"].(string))
				if err != nil {
					return nil, err
				}

				c.Name = p.Args["name"].(string)
				c.Script = p.Args["script"].(string)

				err = c.GenerateTemplateConfigs()
				if err != nil {
					return nil, err
				}

				return c, c.Save()
			},
		},

		"deleteChoreTemplate": &graphql.Field{
			Type: graphql.Boolean,
			Args: idArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := chores.LoadChoreTemplate(p.Args["id"].(string))
				if err != nil {
					return false, err
				}

				err = c.Delete()
				return err == nil, err
			},
		},
	}
}


