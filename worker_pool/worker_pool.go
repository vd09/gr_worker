package worker_pool

type WorkerPool interface {
	AddTask(taskFunc interface{}, params ...interface{}) bool
	AddTaskIfSpaceAvailable(taskFunc interface{}, params ...interface{}) bool
	IsWorkerPoolStopped() bool
	Stop()
	WaitAndStop()
}
