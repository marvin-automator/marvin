package internal

import (
	"fmt"
	"os"
)

func ErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}