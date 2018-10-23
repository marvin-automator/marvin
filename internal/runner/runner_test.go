package runner

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Test that when there are no tasks running, and the context is cancelled, start() returns.
func TestTaskRunner_Run_No_Tasks(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Millisecond)
	tr := NewRunner()
	tr.start(ctx)
}

func TestTaskRunner_Run(t *testing.T) {
	r := require.New(t)
	res := "not run"

	ctx, cancel := context.WithCancel(context.Background())
	tr := NewRunner()
	tr.Run(func(doneCh chan<- struct{}) {
		time.AfterFunc(5*time.Millisecond, func() {
			res = "run"
			doneCh <- struct{}{}
		})
	})

	time.AfterFunc(1*time.Millisecond, cancel)
	tr.start(ctx)

	r.Equal("run", res)
}
