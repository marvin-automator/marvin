package main

import (
	"github.com/bigblind/marvin/handlers"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	port := envy.Get("PORT", "3000")
	app := handlers.App()
	log.Fatal(app.Start(port))
}
