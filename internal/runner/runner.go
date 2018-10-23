package runner

import (
	"context"
	"sync"
)

type TaskRunner struct {
	wg sync.WaitGroup
	doneCh chan struct{}
}

type Task func(doneCh chan<- struct{})

func NewRunner() *TaskRunner {
	return &TaskRunner{doneCh: make(chan struct{}, 10)}
}

func (tr *TaskRunner) start(ctx context.Context) {
	go func() {
		for {
			<- tr.doneCh
			tr.wg.Done()
		}
	}()

	<-ctx.Done()
	tr.wg.Wait()
}

func (tr *TaskRunner) Run(t Task) {
	tr.wg.Add(1)
	t(tr.doneCh)
}