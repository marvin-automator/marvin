package providers

// This file imports all built-in providers, so they get registered.

import (
	_ "github.com/marvin-automator/marvin/actions/providers/web" // register the URL provider
	_ "github.com/marvin-automator/marvin/actions/providers/github" // register the GitHub provider
)
