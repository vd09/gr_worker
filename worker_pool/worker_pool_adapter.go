package worker_pool

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/domain"
	"github.com/vd09/gr_worker/logger"
	"github.com/vd09/gr_worker/worker"
)

const (
	DefaultMinWorkers     = -1
	DefaultMaxWorkers     = 1
	DefaultMaxTasks       = 1
	DefaultIdleTimeout    = 5 * time.Second
	DefaultWorkerStrategy = worker.STANDARD_WORKER
)

type WorkerPoolAdapter struct {
	// context settings
	ctx       context.Context
	cancelCtx context.CancelFunc
	logger    logger.Logger

	// Atomic counters, should be placed first so alignment is guaranteed for atomic operations.
	activeWorkerCount atomic.Int32
	idleWorkerCount   atomic.Int32
	stopped           atomic.Bool

	// Configurable settings
	minWorkers  int32
	maxWorkers  int32
	maxTasks    int32
	idleTimeout time.Duration
	strategy    worker.WorkerStrategy

	// Private properties
	mutex sync.Mutex
	tasks gr_variable.GrChannel[*gr_worker.Task]
}

func NewWorkerPool(options ...Option) (WorkerPool, error) {
	return NewWorkerPoolAdapter(options...)
}

func NewWorkerPoolAdapter(options ...Option) (*WorkerPoolAdapter, error) {
	wp := &WorkerPoolAdapter{
		minWorkers:  DefaultMinWorkers,
		maxWorkers:  DefaultMaxWorkers,
		maxTasks:    DefaultMaxTasks,
		idleTimeout: DefaultIdleTimeout,
		strategy:    DefaultWorkerStrategy,
	}
	wp.stopped.Store(false)
	wp.activeWorkerCount.Store(0)
	wp.idleWorkerCount.Store(0)

	// Apply all options
	for _, opt := range options {
		opt(wp)
	}

	if err := wp.ValidateWorkerPool(); err != nil {
		return nil, err
	}

	if wp.ctx == nil {
		WithContext(context.Background())(wp)
	}
	if wp.minWorkers < 0 {
		wp.minWorkers = wp.maxWorkers
	}
	wp.tasks = gr_variable.NewGrChannelWithLength[*gr_worker.Task](int(wp.maxTasks))

	for i := int32(0); i < wp.minWorkers; i++ {
		wp.startNewWorkerIfRequired()
	}
	return wp, nil
}

func (wp *WorkerPoolAdapter) IsWorkerPoolStopped() bool {
	return wp.stopped.Load()
}

func (wp *WorkerPoolAdapter) Stop() {
	wp.cancelCtx()
	wp.stopped.Store(true)
	wp.tasks.StopWriting()
}

func (wp *WorkerPoolAdapter) WaitAndStop() {
	wp.tasks.StopWriting()
	for {
		if wp.activeWorkerCount.Load() <= 0 {
			wp.cancelCtx()
			wp.stopped.Store(true)
			return
		}
	}
}

func (wp *WorkerPoolAdapter) AddTaskIfSpaceAvailable(taskFunc interface{}, params ...interface{}) bool {
	if wp.IsWorkerPoolStopped() {
		return false
	}

	originalTask := gr_worker.NewTask(taskFunc, params...)
	newTask := gr_worker.NewTask(wp.addConcurrencyDetailsToNewTask, originalTask)

	wp.startNewWorkerIfRequired()
	return wp.tasks.WriteValue(newTask)
}

func (wp *WorkerPoolAdapter) AddTask(taskFunc interface{}, params ...interface{}) bool {
	if wp.IsWorkerPoolStopped() {
		return false
	}

	originalTask := gr_worker.NewTask(taskFunc, params...)
	newTask := gr_worker.NewTask(wp.addConcurrencyDetailsToNewTask, originalTask)

	wp.startNewWorkerIfRequired()
	wp.tasks.MustWriteValue(newTask)
	return true
}

func (wp *WorkerPoolAdapter) addConcurrencyDetailsToNewTask(originalTask *gr_worker.Task) error {
	wp.idleWorkerCount.Add(-1)
	err := originalTask.ExecuteTask()
	wp.idleWorkerCount.Add(1)
	return err
}

func (wp *WorkerPoolAdapter) startNewWorkerIfRequired() {
	if wp.increaseWorkerCount() {
		newWorker := wp.createNewWorker()
		go newWorker.Start()
	}
}

func (wp *WorkerPoolAdapter) createNewWorker() worker.Worker {
	switch wp.strategy {
	case worker.STANDARD_WORKER:
		return worker.NewStandardWorker(wp.ctx, wp.tasks, wp.logger, wp.decreaseWorkerCount)
	case worker.IDEAL_WORKER_TIMEOUT:
		return worker.NewIdealTimeoutWorker(wp.ctx, wp.idleTimeout, wp.logger, wp.tasks, wp.decreaseWorkerCount)
	case worker.SINGLE_TASK_WORKER:
		return worker.NewSingleTaskWorker(wp.ctx, wp.tasks, wp.logger, wp.decreaseWorkerCount)
	}
	return nil
}

//func (wp *WorkerPoolAdapter) canStopWorker(isCtxDone bool) bool {
//	return isCtxDone || wp.decreaseWorkerCount()
//}

func (wp *WorkerPoolAdapter) increaseWorkerCount() bool {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	if wp.IsWorkerPoolStopped() {
		return false
	}
	if wp.activeWorkerCount.Load() >= wp.maxWorkers {
		return false
	}
	if int(wp.idleWorkerCount.Load()) > len(wp.tasks.Receive()) && wp.activeWorkerCount.Load() >= wp.minWorkers {
		return false
	}

	wp.idleWorkerCount.Add(1)
	wp.activeWorkerCount.Add(1)
	return true
}

func (wp *WorkerPoolAdapter) decreaseWorkerCount(workerStatus domain.WorkerStatus) bool {
	wp.mutex.Lock()
	defer wp.mutex.Unlock()

	switch workerStatus {
	case domain.TIMEOUT:
		if wp.idleWorkerCount.Load() <= 0 || (wp.activeWorkerCount.Load() <= wp.minWorkers) {
			return false
		}
	}

	wp.idleWorkerCount.Add(-1)
	wp.activeWorkerCount.Add(-1)
	return true
}
