package runner

import (
	"context"
	"sync"
)

type taskRunner struct {
	wg     sync.WaitGroup
	doneCh chan struct{}
}

type Task func(doneCh chan<- struct{})

type TaskRunner interface {
	Start(ctx context.Context)
	Run(t Task) error
}

func NewRunner() TaskRunner {
	return &taskRunner{doneCh: make(chan struct{}, 100)}
}

func (tr *taskRunner) Start(ctx context.Context) {
	go func() {
		for {
			<-tr.doneCh
			tr.wg.Done()
		}
	}()

	<-ctx.Done()
	tr.wg.Wait()
}

func (tr *taskRunner) Run(t Task) (err error) {
	defer func(){
		v := recover()
		if v != nil {
			tr.doneCh<- struct{}{}
		}
	}()
	tr.wg.Add(1)
	t(tr.doneCh)
	return nil
}
