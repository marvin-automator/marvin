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
	tr.Start(ctx)
}

// Test that Start() returns upon cancelling the context, but waits for running functions to finish.
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
	tr.Start(ctx)

	r.Equal("run", res)
}

// Test that Start() still returns normally when a function passed to Run() panics
func TestTaskRunner_Run_Panic(t *testing.T) {
	r := require.New(t)

	ctx, cancel := context.WithCancel(context.Background())
	tr := NewRunner()
	tr.Run(func(doneCh chan<- struct{}) {
		time.AfterFunc(5*time.Millisecond, func() {
			panic("heyo!")
		})
	})

	time.AfterFunc(1*time.Millisecond, cancel)
	tr.Start(ctx)
}