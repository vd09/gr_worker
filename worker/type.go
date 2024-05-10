package worker

type WorkerStrategy int

const (
	IDEAL_WORKER_TIMEOUT WorkerStrategy = iota
	SINGLE_TASK_WORKER
	STANDARD_WORKER
)
