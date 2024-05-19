package worker

import (
	"context"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/domain"
	"github.com/vd09/gr_worker/logger"
)

type StandardWorker struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    logger.Logger

	isEligibleToStop IsEligibleToStopFunc
	tasks            gr_variable.ReadOnlyGrChannel[*gr_worker.Task]
}

func (sw *StandardWorker) Start() {
	for {
		select {
		case <-sw.ctx.Done():
			if sw.isEligibleToStop(domain.CONTEXT_DONE) {
				return
			}
		case task, ok := <-sw.tasks.Receive():
			switch {
			case ok:
				if err := task.ExecuteTask(); err != nil {
					sw.logger.Printf("[ERROR] Function %#v return non nil result: %#v", task, err)
				}
			case sw.isEligibleToStop(domain.ALL_TASKS_DONE):
				return
			}
		}
	}
}

func NewStandardWorker(parentCtx context.Context, tasks gr_variable.GrChannel[*gr_worker.Task], logger logger.Logger,
	stopFunc IsEligibleToStopFunc) Worker {

	ctx, cancelCtx := context.WithCancel(parentCtx)
	return &StandardWorker{
		ctx:              ctx,
		ctxCancel:        cancelCtx,
		tasks:            tasks,
		logger:           logger,
		isEligibleToStop: stopFunc,
	}
}
