package actions

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeRegistry() *ProviderRegistry {
	bytes := []byte{}

	r := NewRegistry()
	r.AddProvider("provider1", "d", bytes)
	p2 := r.AddProvider("provider2", "pd", bytes)

	p2.AddGroup("g1", "d", bytes)
	g2 := p2.AddGroup("g2", "gd", bytes)

	g2.AddAction("runFoo", "ad", bytes, func(inp struct{}, ctx context.Context) (struct{}, error) {
		return struct{}{}, nil
	})

	return r
}

func TestProviderRegistry_Providers(t *testing.T) {
	reg := makeRegistry()

	r := require.New(t)
	r.Equal(2, len(reg.Providers()))
}

func TestProviderRegistry_GetAction(t *testing.T) {
	reg := makeRegistry()

	a, err := reg.GetAction("provider2", "g2", "runFoo")

	r := require.New(t)
	r.Equal("ad", a.Info().Description)
	r.NoError(err)
}

func TestGroup_AddManualTrigger(t *testing.T) {
	r := makeRegistry()

	g := r.providers["provider2"].groups["g1"]
	g.AddManualTrigger("myTrigger", "test trigger", []byte{},
		func(s struct{}, ctx context.Context) (<-chan struct{}, error) {
			return nil, nil
		})
}

func TestProviderGroupAndActionPath(t *testing.T) {
	reg := makeRegistry()
	r := require.New(t)

	p := reg.Providers()[1]
	g := p.Groups()[1]
	a := g.Actions()[0]

	r.Equal("provider2", p.Info().Path())
	r.Equal("provider2.g2", g.Info().Path())
	r.Equal("provider2.g2.runFoo", a.Info().Path())
}
