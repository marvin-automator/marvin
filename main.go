package main

import (
	"github.com/bigblind/marvin/app"
	"github.com/bigblind/marvin/storage"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	marvin := app.App()
	log.Fatal(marvin.Start(port))
}
