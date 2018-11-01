package chores

import (
	"bytes"
	"encoding/json"
	"github.com/bigblind/v-eight"
	"github.com/gobuffalo/packr"
	"github.com/marvin-automator/marvin/actions"
	"text/template"
)

type RegisteredTrigger struct {
	Provider string
	Group    string
	Action   string
}

type ConfigInput struct {
	Name        string
	Description string
}

type ChoreTemplateConfig struct {
	Triggers []RegisteredTrigger
	Inputs   []ConfigInput
}

type ChoreTemplate struct {
	Name   string
	Id     []byte
	Script string
	Config ChoreTemplateConfig
}

var bp, bpErr = packr.NewBox("./js").FindString("boilerplate.js.template")
var bpTemplate = template.Must(template.New("js").Parse(bp))

func (ct *ChoreTemplate) combineScriptWithBoilerplate(inputs map[string]string) string {
	if bpErr != nil {
		panic(bpErr) // If everything is set up correctly during build, this shouldn't happen.
	}

	w := bytes.NewBuffer([]byte{})
	bpTemplate.Execute(w, struct{
		Providers []actions.Provider
		Inputs map[string]string
	}{actions.Registry.Providers(), inputs})

	s := w.String() + ct.Script
	return s
}

func (ct *ChoreTemplate) GetChoreSnapshot(inputs map[string]string) []byte {
	s := v8.CreateSnapshot(ct.combineScriptWithBoilerplate(inputs))
	return s.Export()
}

func (ct *ChoreTemplate) getConfigV8Value(inputs map[string]string) (*v8.Value, error) {
	is := v8.NewIsolate()
	return is.NewContext().Eval(ct.combineScriptWithBoilerplate(inputs)+";marvin;", "marvin.js")

}

func (ct *ChoreTemplate) GenerateTemplateConfigs() error {
	res, err := ct.getConfigV8Value(nil)
	if err != nil {
		return err
	}

	ct.Config = ChoreTemplateConfig{}

	err = decodeFieldValue(res, "_triggers", &(ct.Config.Triggers))
	if err != nil {
		return err
	}

	return decodeFieldValue(res, "_inputs", &(ct.Config.Inputs))
}

func (ct *ChoreTemplate) GenerateChoreConfig(inputValues map[string]string) (*choreConfig, error) {
	res, err := ct.getConfigV8Value(inputValues)
	if err != nil {
		return nil, err
	}

	cc := choreConfig{}
	cc.Inputs = inputValues

	err = decodeFieldValue(res, "_triggers", &(cc.Triggers))
	return &cc, err
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
