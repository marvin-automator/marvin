package graphql

import (
	"github.com/graphql-go/graphql"
	"github.com/marvin-automator/marvin/internal/chores"
	"time"
)

var (
	idArgs graphql.FieldConfigArgument
	choreType graphql.Output
	choreTemplateType graphql.Output
)
func init() {
	RegisterTypeTransformer(func(t time.Time) string {
		return t.Format(time.RFC3339)
	})


	idArgs = graphql.FieldConfigArgument{
		"id": &graphql.ArgumentConfig{
			Type: graphql.String,
		},
	}

	choreType = CreateOutputTypeFromStruct(chores.Chore{})
	choreTemplateType = CreateOutputTypeFromStruct(chores.ChoreTemplate{})
}

type inputValue struct {
	Name string `json:"name"`
	Value string `json:"value"`
}

func getChoreQueryFields() graphql.Fields {

	choreSettingsType := CreateOutputTypeFromStruct(choreTemplateType).(*graphql.Object)

	choreSettingsType.AddFieldConfig("inputs", &graphql.Field{
		Type: graphql.NewList(CreateOutputTypeFromStruct(inputValue{})),
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

	choreType.(*graphql.Object).AddFieldConfig("logs", &graphql.Field{
		Type: graphql.NewList(CreateOutputTypeFromStruct(chores.ChoreLog{})),
		Args: graphql.FieldConfigArgument{
			"upTo": &graphql.ArgumentConfig{Type: graphql.String},
			"count": &graphql.ArgumentConfig{Type: graphql.Int},
		},
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			n := 0
			if newn, ok := p.Args["count"].(int); ok {
				n = newn
			}

			t := time.Now()
			var err error
			if upTo, ok := p.Args["upTo"].(string); ok {
				t, err = time.Parse(time.RFC3339, upTo)
				if err != nil {
					return nil, err
				}
			}

			chore := p.Source.(*chores.Chore)
			return chore.GetLogsUpTo(t, n)
		},
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


