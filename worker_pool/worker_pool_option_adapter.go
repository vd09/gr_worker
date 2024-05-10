package worker_pool

import (
	"context"
	"time"

	"github.com/vd09/gr_worker/worker"
)

type Option func(*WorkerPoolAdapter)

// MinWorkers allows to change the minimum number of workers of a worker pool
func WithMinWorkers(minWorkers int32) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.minWorkers = minWorkers
	}
}

// MaxWorkers allows to change the minimum number of workers of a worker pool
func WithMaxWorkers(maxWorkers int32) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.maxWorkers = maxWorkers
	}
}

// MaxTasks allows to change the minimum number of workers of a worker pool
func WithMaxTasks(maxTasks int32) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.maxTasks = maxTasks
	}
}

// IdleTimeout allows to change the idle timeout for a worker pool
func WithIdleTimeout(idleTimeout time.Duration) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.idleTimeout = idleTimeout
	}
}

// Worker strategy allows to change the strategy used to resize the pool
func WithWorkerStrategy(strategy worker.WorkerStrategy) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.strategy = strategy
	}
}

// Context configures a parent context on a worker pool to stop all workers when it is cancelled
func WithContext(parentCtx context.Context) Option {
	return func(wp *WorkerPoolAdapter) {
		wp.ctx, wp.cancelCtx = context.WithCancel(parentCtx)
	}
}
