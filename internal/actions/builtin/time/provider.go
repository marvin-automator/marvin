package time

import (
	"context"
	"fmt"
	"github.com/gorhill/cronexpr"
	"github.com/marvin-automator/marvin/actions"
	"time"
)

func init() {
	p := actions.Registry.AddProvider("time", "Time-related actions", []byte{})
	g := p.AddGroup("cron", "Cron-related tasks", []byte{})
	g.AddManualTrigger("onCronSchedule", "Schedule a function to run on an interval based on a Cron expression.", []byte{}, cronTrigger)
	g.AddAction("nextScheduledTime", "Get the next time a cron trigger with the given expression would run.", []byte{}, getNextExecutionTime)
}

func cronTrigger(exprs string, ctx context.Context) (<-chan time.Time, error) {
	expr, err := cronexpr.Parse(exprs)
	if err != nil {
		return nil, err
	}

	out := make(chan time.Time, 10)
	var f func()
	f = func() {
		now := time.Now()
		n := expr.Next(now)
		duration := n.Sub(now)
		fmt.Printf("Next scheduled event in %v\n", duration)
		t := time.NewTimer(duration)

		select {
		case <-ctx.Done():
		case now = <-t.C:
			fmt.Println("Sending event")
			out <- now
			fmt.Println("Event sent")
			go f()
		}
	}

	go f()
	return out, nil
}

func getNextExecutionTime(exprs string, ctx context.Context) (time.Time, error) {
	expr, err := cronexpr.Parse(exprs)
	if err != nil {
		return time.Time{}, err
	}

	return expr.Next(time.Now()), nil
}
