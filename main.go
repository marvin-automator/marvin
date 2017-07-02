package main

import (
	"github.com/bigblind/marvin/actions"
	"github.com/gobuffalo/envy"
	"log"
)

func main() {
	port := envy.Get("PORT", "3000")
	app := actions.App()
	log.Fatal(app.Start(port))
}
