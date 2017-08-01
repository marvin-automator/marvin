package main

import (
	"github.com/marvin-automator/marvin/app"
	"github.com/marvin-automator/marvin/storage"
	"github.com/gobuffalo/envy"
	"log"
	"github.com/marvin-automator/marvin/actions/interactors/execution"
	"github.com/marvin-automator/marvin/app/domain"

	// import the built-in providers so they get registered
	_ "github.com/marvin-automator/marvin/actions/providers"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	marvin := app.App()
	execution.SetupExecutionEnvironment(marvin.Context, domain.LoggerFromBuffaloLogger(marvin.Logger))
	log.Fatal(marvin.Start(port))
}
