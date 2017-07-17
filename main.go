package main

import (
	"github.com/bigblind/marvin/app"
	"github.com/bigblind/marvin/storage"
	"github.com/gobuffalo/envy"
	"log"
	"github.com/bigblind/marvin/actions/interactors/execution"
	"github.com/bigblind/marvin/app/domain"

	// import the built-in providers so they get registered
	_ "github.com/bigblind/marvin/actions/providers"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	marvin := app.App()
	execution.SetupExecutionEnvironment(marvin.Context, domain.LoggerFromBuffaloLogger(marvin.Logger))
	log.Fatal(marvin.Start(port))
}
