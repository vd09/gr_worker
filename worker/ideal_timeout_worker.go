package worker

import (
	"context"
	"time"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/logger"
)

type IdealTimeoutWorker struct {
	ctx       context.Context
	ctxCancel context.CancelFunc

	logger       logger.Logger
	idealTimeout time.Duration
	timer        *time.Timer

	isEligibleToStop IsEligibleToStopFunc

	tasks gr_variable.ReadOnlyGrChannel[*gr_worker.Task]
	errs  gr_variable.WriteOnlyGrChannel[error]
}

func (itw *IdealTimeoutWorker) Start() {
	itw.timer = time.NewTimer(itw.idealTimeout)
	for {
		select {
		case <-itw.ctx.Done():
			if itw.isEligibleToStop(true) {
				return
			}
		case <-itw.timer.C:
			if itw.isEligibleToStop(false) {
				return
			}
		case task, ok := <-itw.tasks.Receive():
			if !ok {
				return
			}
			if err := task.ExecuteTask(); err != nil {
				itw.logger.Printf("[ERROR] Function %#v return non nil result: %#v", task, err)
			}
			itw.timer = time.NewTimer(itw.idealTimeout)
		}
	}
}

func NewIdealTimeoutWorker(parentCtx context.Context, idleTimeout time.Duration, logger logger.Logger,
	tasks gr_variable.GrChannel[*gr_worker.Task], stopFunc IsEligibleToStopFunc) Worker {

	ctx, cancelCtx := context.WithCancel(parentCtx)
	return &IdealTimeoutWorker{
		ctx:              ctx,
		ctxCancel:        cancelCtx,
		tasks:            tasks,
		logger:           logger,
		idealTimeout:     idleTimeout,
		isEligibleToStop: stopFunc,
	}
}
