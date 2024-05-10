package worker_pool

import "errors"

var (
	ErrMaxWorkers  = errors.New("max workers can't be less than one")
	ErrMinWorkers  = errors.New("min workers is greater than max workers")
	ErrMaxTasks    = errors.New("max tasks can't be less than one")
	ErrIdleTimeout = errors.New("max tasks can't be less than zero")
)

func (wp *WorkerPoolAdapter) validateMaxWorkers() error {
	if wp.maxWorkers <= 0 {
		return ErrMaxWorkers
	}
	return nil
}

func (wp *WorkerPoolAdapter) validateMinWorkers() error {
	if wp.minWorkers > wp.maxWorkers {
		return ErrMinWorkers
	}
	return nil
}

func (wp *WorkerPoolAdapter) validateMaxTasks() error {
	if wp.maxTasks <= 0 {
		return ErrMaxTasks
	}
	return nil
}

func (wp *WorkerPoolAdapter) validateIdleTimeout() error {
	if wp.idleTimeout < 0 {
		return ErrIdleTimeout
	}
	return nil
}

func (wp *WorkerPoolAdapter) ValidateWorkerPool() error {
	if err := wp.validateMaxWorkers(); err != nil {
		return err
	}
	if err := wp.validateMinWorkers(); err != nil {
		return err
	}
	if err := wp.validateMaxTasks(); err != nil {
		return err
	}
	if err := wp.validateIdleTimeout(); err != nil {
		return err
	}
	return nil
}
