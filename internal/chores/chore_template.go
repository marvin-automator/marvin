package chores

import (
	"bytes"
	"encoding/json"
	"github.com/joesonw/js8"
	"github.com/gobuffalo/packr"
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/internal"
	"github.com/marvin-automator/marvin/internal/db"
	"github.com/pkg/errors"
	"text/template"
	"time"
)

// RegisteredTrigger holds information about a trigger, registered in a template.
type RegisteredTrigger struct {
	Provider string `json:"provider"`
	Group    string `json:"group"`
	Action   string `json:"action"`
}

// ConfigInput describes a parameter that is configurable by the user when creating a chore from a template.
type ConfigInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ChoreTemplateConfig holds configuration information about a trigger.
type ChoreTemplateConfig struct {
	Triggers []RegisteredTrigger `json:"triggers"`
	Inputs   []ConfigInput       `json:"inputs"`
}

// ChoreTemplate is a template for a chore (how surprising!)
type ChoreTemplate struct {
	Name    string              `json:"name"`
	Id      string              `json:"id"`
	Created time.Time           `json:"created"`
	Script  string              `json:"script"`
	Config  ChoreTemplateConfig `json:"templateSettings"`
}

var (
	errTemplateNotFound = errors.New("Template not found")
)

// NewChoreTemplate creates a new chore template with the given name and string. It executes the script to determine
// which triggers are registered, and which inputs are requested.
func NewChoreTemplate(name, script string) (*ChoreTemplate, error) {
	id, err := internal.NewId()
	if err != nil {
		return nil, err
	}

	ct := ChoreTemplate{
		Name:    name,
		Script:  script,
		Id:      id,
		Created: time.Now(),
	}

	err = ct.GenerateTemplateConfigs()
	return &ct, err
}

var templateCache = make(map[string]*ChoreTemplate)
var cacheLoaded = false

const templateStoreName = "chore_templates"

// LoadChoreTemplate loads a single template from the database
func LoadChoreTemplate(id string) (*ChoreTemplate, error) {
	ct, ok := templateCache[id]
	if ok {
		return ct, nil
	}

	s := db.GetStore(templateStoreName)
	ct = new(ChoreTemplate)
	err := s.Get(id, ct)
	if _, ok := err.(db.KeyNotFoundError); ok {
		return ct, errTemplateNotFound
	}

	if err != nil {
		return ct, err
	}

	templateCache[id] = ct
	return ct, nil
}

// GetChoreTemplates gets all chore templates from the database.
func GetChoreTemplates() ([]*ChoreTemplate, error) {
	results := make([]*ChoreTemplate, 0, len(templateCache))

	if cacheLoaded {
		for _, ct := range templateCache {
			results = append(results, ct)
		}
		return results, nil
	}

	s := db.GetStore(templateStoreName)
	ct := new(ChoreTemplate)
	err := s.EachKeyWithPrefix("", ct, func(key string) error {
		ctcopy := *ct
		templateCache[key] = &ctcopy
		results = append(results, &ctcopy)
		return nil
	})

	if err != nil {
		return nil, err
	}

	cacheLoaded = true
	return results, nil
}

// Save stores a chore template to the database.
func (ct *ChoreTemplate) Save() error {
	s := db.GetStore("chore_templates")
	templateCache[ct.Id] = ct
	return s.Set(ct.Id, ct)
}

// Delete removes a chore template from the database.
func (ct *ChoreTemplate) Delete() error {
	s := db.GetStore(templateStoreName)
	err := s.Delete(ct.Id)
	if err != nil {
		return err
	}
	delete(templateCache, ct.Id)
	return nil
}

var bp, bpErr = packr.NewBox("./js").FindString("boilerplate.js.template")
var bpTemplate = template.Must(template.New("js").Parse(bp))

// combineScriptWithBoilerplate generates JavaScript code for a chore by combining the chore's script with
// some setup code that makes actions and triggers, as well as marvin.input and other utilities, available to the script.
func (ct *ChoreTemplate) combineScriptWithBoilerplate(inputs map[string]string) string {
	if bpErr != nil {
		panic(bpErr) // If everything is set up correctly during build, this shouldn't happen.
	}

	w := bytes.NewBuffer([]byte{})
	bpTemplate.Execute(w, struct {
		Providers []actions.Provider
		Inputs    map[string]string
	}{actions.Registry.Providers(), inputs})

	s := w.String() + ct.Script
	//s = v8console.WrapForSnapshot(s)
	return s
}

// GetChoreSnapshot generates a v8 Snapshot with all triggers in the template registered.
func (ct *ChoreTemplate) GetChoreSnapshot(inputs map[string]string) []byte {
	s := v8.CreateSnapshot(ct.combineScriptWithBoilerplate(inputs))
	return s.Export()
}

// getConfigV8Value executes the tamplate to get the value needed to create a ChoreTemplateConfig
func (ct *ChoreTemplate) getConfigV8Value(inputs map[string]string) (*v8.Value, error) {
	is := v8.NewIsolate()
	return is.NewContext().Eval(ct.combineScriptWithBoilerplate(inputs)+";marvin;", "marvin.js")

}

// GenerateTemplateConfig runs the script to generate a ChoreTemplateConfig
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

// GenerateChoreConfig runs the template's script with the given inputs to generate a chore config.
func (ct *ChoreTemplate) GenerateChoreConfig(inputValues map[string]string) (*ChoreConfig, error) {
	res, err := ct.getConfigV8Value(inputValues)
	if err != nil {
		return nil, err
	}

	cc := ChoreConfig{}
	cc.Inputs = inputValues

	err = decodeFieldValue(res, "_triggers", &(cc.Triggers))
	return &cc, err
}

// decodeFieldValue is a helper that takes an attribute of a v8 object, and decodes it.
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
