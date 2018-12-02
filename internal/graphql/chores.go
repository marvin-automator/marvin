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

var choreSettingsType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ChoreSettings",
	Fields: graphql.BindFields(chores.ChoreConfig{}),
})

var choreTemplateType = graphql.NewObject(graphql.ObjectConfig{
	Name:        "ChoreTemplate",
	Fields:      graphql.BindFields(chores.ChoreTemplate{}),
	Description: "Describes a template for chores.",
})

type inputValue struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func getChoreQueryFields() graphql.Fields {

	choreType.AddFieldConfig("template", &graphql.Field{
		Type: choreTemplateType,
	})

	choreSettingsType.AddFieldConfig("inputs", &graphql.Field{
		Type: graphql.NewList(graphql.NewObject(graphql.ObjectConfig{
			Name: "InputValue",
			Fields: graphql.BindFields(inputValue{}),
		})),
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			cf := p.Source.(chores.ChoreConfig)
			res := make([]inputValue, len(cf.Inputs))
			i := 0
			for name, value := range cf.Inputs {
				res[i] = inputValue{name, value}
				i += 1
			}

			return res, nil
		},
	})

	choreType.AddFieldConfig("choreSettings", &graphql.Field{
		Type: choreSettingsType,
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

		// chore mutations
		"createChore": &graphql.Field{
			Type: choreType,
			Args: graphql.FieldConfigArgument{
				"templateId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				"inputs": &graphql.ArgumentConfig{Type: graphql.NewList(graphql.NewInputObject(graphql.InputObjectConfig{
					Name: "ChoreInput",
					Description: "Defines a value for a particular chore input defined in a template",
					Fields: graphql.InputObjectConfigFieldMap{
						"name": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
						"value": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
					},
				}))},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				ctId := p.Args["templateId"].(string)
				ct, err := chores.LoadChoreTemplate(ctId)
				if err != nil {
					return nil, err
				}

				inputs := make(map[string]string)
				gqlInputs := p.Args["inputs"].([]interface{})
				for _, i := range gqlInputs {
					im := i.(map[string]interface{})
					inputs[im["name"].(string)] = im["value"].(string)
				}

				c, err := chores.FromTemplate(ct, p.Args["name"].(string), inputs)
				if err != nil {
					return nil, err
				}

				return c, c.Save()
			},
		},

		"deleteChore": &graphql.Field{
			Type: graphql.Boolean,
			Args: idArgs,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				c, err := chores.GetChore(p.Args["id"].(string))
				if err != nil {
					return false, err
				}

				err = c.Delete()
				return err == nil, err
			},
		},
	}
}


