package handlers_test

import (
	"testing"

	"github.com/bigblind/marvin/handlers"
	"github.com/gobuffalo/suite"
)

type ActionSuite struct {
	*suite.Action
}

func Test_ActionSuite(t *testing.T) {
	as := &ActionSuite{suite.NewAction(handlers.App())}
	suite.Run(t, as)
}
