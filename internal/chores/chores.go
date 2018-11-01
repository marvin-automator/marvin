package chores

import (
	"encoding/json"
	"github.com/marvin-automator/marvin/actions"
	"github.com/satori/go.uuid"
	"reflect"
)

type choreTrigger struct {
	RegisteredTrigger
	Input interface{}
}

func (ct *choreTrigger) UnmarshalJSON(data []byte) error {
	var rt RegisteredTrigger
	err := json.Unmarshal(data, &rt)
	if err != nil {
		return err
	}

	ct.RegisteredTrigger = rt

	a, err := actions.Registry.GetAction(rt.Provider, rt.Group, rt.Action)
	if err != nil {
		return err
	}

	inputHolder := struct{Input interface{}}{}
	inputHolder.Input = reflect.New(a.Info().InputType).Interface()
	err = json.Unmarshal(data, &inputHolder)
	if err != nil {
		return err
	}

	ct.Input = inputHolder.Input

	return nil
}

type choreConfig struct {
	Inputs map[string]string
	Triggers []choreTrigger
}

type Chore struct {
	Name string
	Id []byte
	Template ChoreTemplate
	Config choreConfig
	Snapshot []byte
}

func FromTemplate(ct *ChoreTemplate, name string, inputs map[string]string) (*Chore, error) {
	conf, err := ct.GenerateChoreConfig(inputs)
	if err != nil {
		return nil, err
	}

	return &Chore{
		Name: name,
		Config: *conf,
		Template: *ct,
		Snapshot: ct.GetChoreSnapshot(inputs),
		Id: uuid.NewV4().Bytes(),
	}, nil
}


