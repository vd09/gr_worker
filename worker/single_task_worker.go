package worker

import (
	"context"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/domain"
	"github.com/vd09/gr_worker/logger"
)

type SingleTaskWorker struct {
	ctx       context.Context
	ctxCancel context.CancelFunc
	logger    logger.Logger

	isEligibleToStop IsEligibleToStopFunc
	tasks            gr_variable.ReadOnlyGrChannel[*gr_worker.Task]
}

func (btw *SingleTaskWorker) Start() {
	assignedTask, ok := btw.tasks.ReadValue()
	if !ok {
		return
	}

	for {
		select {
		case <-btw.ctx.Done():
			if btw.isEligibleToStop(domain.CONTEXT_DONE) {
				return
			}
		default:
			if err := assignedTask.ExecuteTask(); err != nil {
				btw.logger.Printf("[ERROR] Function %#v return non nil result: %#v", assignedTask, err)
			}
		}
	}
}

func NewSingleTaskWorker(parentCtx context.Context, tasks gr_variable.GrChannel[*gr_worker.Task], logger logger.Logger,
	stopFunc IsEligibleToStopFunc) Worker {

	ctx, cancelCtx := context.WithCancel(parentCtx)
	return &SingleTaskWorker{
		ctx:              ctx,
		ctxCancel:        cancelCtx,
		tasks:            tasks,
		logger:           logger,
		isEligibleToStop: stopFunc,
	}
}
