package main

import (
	"github.com/gobuffalo/envy"
	"log"
	"github.com/bigblind/marvin/storage"
)

func main() {
	port := envy.Get("PORT", "3000")
	storage.Setup()
	App()
	log.Printf("App: %#v", app)
	log.Fatal(app.Start(port))
}
