package worker_pool

import (
	"testing"
	"time"
)

func TestWorkerPoolAdapter_ValidateWorkerPool(t *testing.T) {
	tests := []struct {
		name          string
		minWorkers    int32
		maxWorkers    int32
		maxTasks      int32
		idleTimeout   time.Duration
		expectedError error
	}{
		{"ValidOptions", 1, 2, 3, 5 * time.Second, nil},
		{"InvalidMaxWorkers", 0, 0, 3, 5 * time.Second, ErrMaxWorkers},
		{"InvalidMinWorkers", 2, 1, 3, 5 * time.Second, ErrMinWorkers},
		{"InvalidMaxTasks", 1, 2, 0, 5 * time.Second, ErrMaxTasks},
		{"InvalidIdleTimeout", 1, 2, 3, -1 * time.Second, ErrIdleTimeout},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			wp := &WorkerPoolAdapter{
				minWorkers:  test.minWorkers,
				maxWorkers:  test.maxWorkers,
				maxTasks:    test.maxTasks,
				idleTimeout: test.idleTimeout,
			}
			err := wp.ValidateWorkerPool()
			if err != test.expectedError {
				t.Errorf("unexpected error, got: %v, want: %v", err, test.expectedError)
			}
		})
	}
}
