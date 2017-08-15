package domain

import "github.com/urakozz/go-json-schema-generator"

type ActionSchemas struct {
	InputSchema  string `json:"inputSchema"`
	OutputSchema string `json:"outputSchema"`
	ConfigSchema string `json:"configSchema"`
}

func GetActionSchemas(a BaseAction, ac ActionContext) ActionSchemas {
	schemas := ActionSchemas{}

	ot := a.OutputType(ac)
	schemas.OutputSchema = generator.Generate(ot)

	if a.Meta().IsTrigger {
		ta := a.(Trigger)
		ct := ta.ConfigType(ac)
		schemas.ConfigSchema = generator.Generate(ct)
	} else {
		aa := a.(Action)
		it := aa.InputType(ac)
		schemas.InputSchema = generator.Generate(it)
	}

	return schemas
}
