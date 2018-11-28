package internal

import (
	"fmt"
	"os"
	"time"
)

func ErrorAndExit(err error) {
	if err != nil {
		fmt.Println(err)
		// Give the error time to be printed. (I don't know why this is necessary, but sometimes, the error doesn't get printed without this.
		c := time.After(time.Millisecond*50)
		<- c
		os.Exit(1)
	}
}
