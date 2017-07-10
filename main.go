package main

import (
	"github.com/gobuffalo/envy"
	"log"
	"github.com/bigblind/marvin/storage"
	"github.com/bigblind/marvin/app"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	marvin := app.App()
	log.Fatal(marvin.Start(port))
}
