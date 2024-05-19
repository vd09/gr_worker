package domain

type WorkerStatus string

const (
	CONTEXT_DONE   WorkerStatus = "CONTEXT_DONE"
	ALL_TASKS_DONE              = "ALL_TASKS_DONE"
	TIMEOUT                     = "TIMEOUT"
)
