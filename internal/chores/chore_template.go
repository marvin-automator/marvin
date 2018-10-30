package chores

import (
	"encoding/json"
	"github.com/bigblind/v-eight"
	"github.com/gobuffalo/packr"
)

type RegisteredTrigger struct {
	Provider string
	Group string
	Action string
}

type ConfigInput struct {
	Name string
	Description string
}

type ChoreTemplateConfig struct {
	Triggers []RegisteredTrigger
	Inputs []ConfigInput
}


type ChoreTemplate struct {
	Name string
	Id []byte
	Script string
	Config ChoreTemplateConfig
}

var bp, bpErr = packr.NewBox("./js").FindString("boilerplate.js")
func (ct *ChoreTemplate) combineScriptWithBoilerplate() string {
	if bpErr != nil {
		panic(bpErr) // If everything is set up correctly during build, this shouldn't happen.
	}

	return bp + ct.Script
}

func (ct *ChoreTemplate) GenerateConfigs() error {
	is := v8.NewIsolate()
	res, err := is.NewContext().Eval(ct.combineScriptWithBoilerplate() + ";marvin;", "marvin.js")

	if err != nil {
		return err
	}

	ct.Config = ChoreTemplateConfig{}

	err = decodeFieldValue(res, "_triggers", &(ct.Config.Triggers))
	if err != nil {
		return err
	}

	err = decodeFieldValue(res, "_inputs", &(ct.Config.Inputs))
	if err != nil {
		return err
	}

	return nil
}

func decodeFieldValue(value *v8.Value, field string, ptr interface{}) error {
	v, err := value.Get(field)
	if err != nil {
		return err
	}

	data, err := v.MarshalJSON()
	if err != nil {
		return err
	}

	return json.Unmarshal(data, ptr)
}
