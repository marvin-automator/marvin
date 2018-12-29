package chores

import (
	"context"
	"github.com/marvin-automator/marvin/actions"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type TestValue struct {
	S string `json:"s"`
}

var donech = make(chan TestValue)

func addCustomTriggerAndAction() {
	b := []byte{} // We don't care about the svg icon here, so just use an empty slice of bytes

	p := actions.Registry.AddProvider("testProvider", "", b)
	g := p.AddGroup("testGroup", "", b)

	g.AddManualTrigger("onTestTrigger", "", b, func(tv TestValue, ctx context.Context) (<-chan TestValue, error) {
		c := make(chan TestValue)

		go func() {
			c <- TestValue{tv.S + "_triggered"}
		}()

		return c, nil
	})

	g.AddAction("testAction", "", b, func(tv TestValue, ctx context.Context) (TestValue, error) {
		donech <- tv
		return tv, nil
	})

}

func createTemplate() (*ChoreTemplate, error) {
	return NewChoreTemplate("test", `
let inp = marvin.input("myInput", "Just some input")
testProvider.testGroup.onTestTrigger({s: inp}, (r) => {
	testProvider.testGroup.testAction({s: r.s});
})
`)
}

func TestChore_Start(t *testing.T) {
	r := require.New(t)

	addCustomTriggerAndAction()

	ct, err := createTemplate()
	r.NoError(err)

	c, err := FromTemplate(ct, "My Chore", map[string]string{"myInput": "input"})
	r.NoError(err)

	// Basic checks to see that values are set as expected
	r.Equal("My Chore", c.Name)
	r.Equal(map[string]string{"myInput": "input"}, c.Config.Inputs)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*5)
	c.Start(ctx)

	select {
	case <-ctx.Done():
		r.FailNow("The deadline expired, and the action still hasn't been called.")
	case tv := <-donech:
		r.Equal(TestValue{"input_triggered"}, tv)
	}
}
