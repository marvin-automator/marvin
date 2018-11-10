package chores

import (
	"context"
	"github.com/marvin-automator/marvin/actions"
	internalActions "github.com/marvin-automator/marvin/internal/actions"
	"testing"

	"github.com/stretchr/testify/require"
)

var testScript = `
var i1 = marvin.input("my_input", "description");
var i2 = marvin.input("another_input", "description2");

myProvider.myGroup.aTrigger({s: i2}, (e) => e);
`

func setupRegistry() {
	actions.Registry = internalActions.NewRegistry()
	bs := []byte{}
	p := actions.Registry.AddProvider("myProvider", "", bs)
	g := p.AddGroup("myGroup", "", bs)

	g.AddAction("anAction", "", bs, func(s struct{}, ctx context.Context) (struct{}, error) {
		return struct{}{}, nil
	})
	g.AddManualTrigger("aTrigger", "", bs, func(s struct{S string}, ctx context.Context) (<-chan struct{}, error) {
		return nil, nil
	})
}

func TestChoreTemplate_GenerateTemplateConfigs(t *testing.T) {
	setupRegistry()

	ct := &ChoreTemplate{}
	ct.Script = testScript
	err := ct.GenerateTemplateConfigs()

	r := require.New(t)

	r.NoError(err)

	inputs := []ConfigInput{
		{"my_input", "description"},
		{"another_input", "description2"},
	}
	r.Equal(inputs, ct.Config.Inputs)

	triggers := []RegisteredTrigger{
		{"myProvider", "myGroup", "aTrigger"},
	}
	r.Equal(triggers, ct.Config.Triggers)
}


func TestChoreTemplate_GenerateChoreConfig(t *testing.T) {
	setupRegistry()

	ct := &ChoreTemplate{}
	ct.Script = testScript
	cc, err := ct.GenerateChoreConfig(map[string]string{
		"my_input": "val1",
		"another_input": "val2",
	})

	r := require.New(t)

	r.NoError(err)
	triggers := []choreTrigger{
		{RegisteredTrigger{
			"myProvider",
			"myGroup",
			"aTrigger"},
		struct{S string}{"val2"}},
	}
	r.Equal(triggers, cc.Triggers)
}
