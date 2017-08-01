package interactors

import (
	"github.com/marvin-automator/marvin/actions"
	"github.com/marvin-automator/marvin/actions/domain"
	"github.com/stretchr/testify/require"
	"testing"
)

func makeProviderMeta(name, description, key string) domain.ProviderMeta {
	return domain.ProviderMeta{
		Name:        name,
		Description: description,
		Key:         key,
	}
}

func makeActionMeta(name, description, key string, flags bool) domain.ActionMeta {
	return domain.ActionMeta{
		Name:            name,
		Description:     description,
		Key:             key,
		RequiresTestRun: flags,
		IsTrigger:       flags,
	}
}

func TestGetActionGroups(t *testing.T) {
	// Set up a registry and 2 providers
	mr := actions.NewMockRegistry()
	mp1 := actions.NewMockProvider()
	mp2 := actions.NewMockProvider()

	// Set up a list of ActionMetas for groups to return
	acs1 := []domain.ActionMeta{
		makeActionMeta("ac11", "", "ac11", true),
		makeActionMeta("ac12", "", "ac12", false),
	}
	acs2 := []domain.ActionMeta{
		makeActionMeta("ac21", "", "ac21", true),
		makeActionMeta("ac22", "", "ac22", false),
	}
	acs3 := []domain.ActionMeta{
		makeActionMeta("ac31", "", "ac31", true),
		makeActionMeta("ac32", "", "ac32", false),
	}
	acs4 := []domain.ActionMeta{
		makeActionMeta("ac41", "", "ac41", true),
	}

	// Set up groups for the providers to return
	g11 := actions.NewMockGroup()
	g11.On("Actions").Return(acs1)
	g12 := actions.NewMockGroup()
	g12.On("Actions").Return(acs2)
	g21 := actions.NewMockGroup()
	g21.On("Actions").Return(acs3)
	g22 := actions.NewMockGroup()
	g22.On("Actions").Return(acs4)
	g11.On("Name").Return("g11")
	g12.On("Name").Return("g12")
	g21.On("Name").Return("g21")
	g22.On("Name").Return("g22")

	// Make the providers return the groups we just set up
	mp1.On("Groups").Return([]domain.Group{g11, g12})
	mp2.On("Groups").Return([]domain.Group{g21, g22})

	// Make the registry return the providers
	mr.On("Providers").Return([]domain.ProviderMeta{
		makeProviderMeta("", "", "p1"),
		makeProviderMeta("", "", "p2"),
	})
	mr.On("Provider", "p1").Return(mp1)
	mr.On("Provider", "p2").Return(mp2)

	// Now time for claling the function we're actually testing
	interactor := Registry{mr}
	gs := interactor.GetActionGroups()

	// And asserting that it returns all the groups
	require.Equal(t, []Group{
		Group{"g11", "p1", acs1},
		Group{"g12", "p1", acs2},
		Group{"g21", "p2", acs3},
		Group{"g22", "p2", acs4},
	}, gs)
}
