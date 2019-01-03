package icons

import "github.com/gobuffalo/packr"

var box = packr.NewBox(".")

func Get(name string) []byte {
	bytes, _ := box.Find(name)
	return bytes
}
