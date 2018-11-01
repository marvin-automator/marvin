package chores

import (
	"context"
	"github.com/marvin-automator/marvin/actions"
	internal_actions "github.com/marvin-automator/marvin/internal/actions"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChoreTemplate_GetConfigs(t *testing.T) {
	actions.Registry = internal_actions.NewRegistry()
	bs := []byte{}
	p := actions.Registry.AddProvider("myProvider", "", bs)
	g := p.AddGroup("myGroup", "", bs)

	g.AddAction("anAction", "", bs, func(s struct{}, ctx context.Context) (struct{}, error) {
		return struct{}{}, nil
	})
	g.AddManualTrigger("aTrigger", "", bs, func(s struct{}, ctx context.Context) (<-chan struct{}, error) {
		return nil, nil
	})

	ct := &ChoreTemplate{}
	ct.Script = `
var i1 = marvin.input("my_input", "description");
var i2 = marvin.input("another_input", "description2");

myProvider.myGroup.aTrigger({}, (e) => e);

`
	err := ct.GenerateConfigs()

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
