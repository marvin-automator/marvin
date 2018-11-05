package chores

import (
	"context"
	"fmt"
	"github.com/marvin-automator/marvin/internal/config"
	"github.com/radovskyb/watcher"
	"time"
)

func StartTemplateWatcher(ctx context.Context) {
	
}

func startWatcher(ctx context.Context) {
	w := watcher.New()
	w.Add(config.TemplateDir)

	w.Start(time.Second*15)

	for {
		select {
		case <-ctx.Done():
			w.Close()
			return
		case e := <-w.Event:
			handleEvent(e)
		case err := <-w.Error:
			fmt.Println(err)
		}
	}
}

func handleEvent(e watcher.Event) {

}