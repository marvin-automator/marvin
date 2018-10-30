package chores

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestChoreTemplate_GetConfigs(t *testing.T) {
	ct := &ChoreTemplate{}
	ct.Script = `
marvin.input("my_input", "description");
marvin.input("another_input", "description2");
`
	err := ct.GenerateConfigs()

	r :=  require.New(t)

	r.NoError(err)
	inputs := []ConfigInput{
		ConfigInput{"my_input", "description"},
		ConfigInput{"another_input", "description2"},
	}
	r.Equal(inputs, ct.Config.Inputs)
}