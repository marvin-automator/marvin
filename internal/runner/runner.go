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
	Run(t Task)
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

func (tr *taskRunner) Run(t Task) {
	tr.wg.Add(1)
	t(tr.doneCh)
}

var runner TaskRunner
var cancelRunner context.CancelFunc

func SetRunner(r TaskRunner, ctx context.Context) {
	if cancelRunner != nil {
		cancelRunner()
	}

	runner = r
	ctx, cancelRunner = context.WithCancel(ctx)
	go r.Start(ctx)
}

func Run(t Task) {
	runner.Run(t)
}

func Stop() {
	cancelRunner()
}

func init() {
	SetRunner(NewRunner(), context.Background())
}
