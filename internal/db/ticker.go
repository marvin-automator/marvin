package db

import (
	"context"
	"time"
)

func PeriodicallyRunGC(ctx context.Context) {
	t := time.NewTicker(time.Minute)
	for {
		select {
		case <-t.C:
			println("Running DB GC cycle.")
			doGC()
		case <-ctx.Done():
			return
		}
	}
}

func doGC() {
	if db == nil {
		return
	}

	var err error
	for err == nil {
		err = db.RunValueLogGC(0.5)
	}
}
