package main

import (
	_ "github.com/marvin-automator/marvin/internal/actions" // make sure the registry gets initialized.

	_ "github.com/marvin-automator/marvin/internal/actions/builtin/http"
	_ "github.com/marvin-automator/marvin/internal/actions/builtin/time"
	_ "github.com/marvin-automator/marvin/internal/chores"
)
