package main

import (
	"github.com/bigblind/marvin/app"
	"github.com/bigblind/marvin/storage"
	"github.com/gobuffalo/envy"
	"log"
	"github.com/bigblind/marvin/actions/interactors/execution"
	"github.com/bigblind/marvin/app/interactors"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	marvin := app.App()
	execution.SetupExecutionEnvironment(marvin.Context, interactors.LoggerFromBuffaloLogger(marvin.Logger))
	log.Fatal(marvin.Start(port))
}
