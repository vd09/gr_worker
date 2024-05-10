package worker_pool

import (
	"fmt"
	"testing"
	"time"
)

func TestWorkerPoolAdapter_AddTask_IsTaskAdded(t *testing.T) {
	// Create a worker pool with default options
	wp, err := NewWorkerPool(
		WithMaxWorkers(3),
		WithMaxTasks(2),
	)
	if err != nil {
		t.Fatalf("error creating worker pool: %v", err)
	}

	defer wp.Stop()
	taskAdded := wp.AddTask(func() {})
	if !taskAdded {
		t.Error("task not added to worker pool")
	}
}

func TestWorkerPoolAdapter_Stop(t *testing.T) {
	// Create a worker pool with default options
	wp, err := NewWorkerPool()
	if err != nil {
		t.Fatalf("error creating worker pool: %v", err)
	}

	// Stop the worker pool
	wp.Stop()

	// Verify that the worker pool is stopped
	if !wp.IsWorkerPoolStopped() {
		t.Error("worker pool not stopped")
	}
}

func TestWorkerPoolAdapter_AddTaskIfSpaceAvailable(t *testing.T) {
	wp, err := NewWorkerPool(
		WithMaxWorkers(1),
		WithMaxTasks(1),
	)
	if err != nil {
		t.Fatalf("error creating worker pool: %v", err)
	}

	defer wp.Stop()
	basicFunction := func() { time.Sleep(9 * time.Second) }
	taskAdded := wp.AddTaskIfSpaceAvailable(basicFunction)
	if !taskAdded {
		t.Error("task not added to worker pool")
	}

	taskAdded = wp.AddTaskIfSpaceAvailable(basicFunction)
	if taskAdded {
		t.Error("task added to worker pool but space doesn't exist")
	}
}

func TestWorkerPoolAdapter_AddTask(t *testing.T) {
	// Create a worker pool with default options
	wp, err := NewWorkerPoolAdapter(
		WithMinWorkers(1),
		WithMaxWorkers(3),
		WithMaxTasks(2),
	)
	if err != nil {
		t.Fatalf("error creating worker pool: %v", err)
	}

	defer wp.Stop()

	basicFunction := func() { time.Sleep(20 * time.Second) }

	for i := 0; i < 3; i++ {
		taskAdded := wp.AddTask(basicFunction)
		time.Sleep(time.Millisecond)
		if !taskAdded {
			t.Error("task not added to worker pool")
		}
		if wp.activeWorkerCount.Load() != int32(i+1) {
			t.Error(fmt.Sprintf("total worker (%d) are not matching as expected (%d)", wp.activeWorkerCount.Load(), i+1))
		}
	}
	taskAdded := wp.AddTask(basicFunction)
	if !taskAdded {
		t.Error("task not added to worker pool")
	}
	if wp.activeWorkerCount.Load() != 3 {
		t.Error(fmt.Sprintf("total worker (%d) are not matching as expected (3)", wp.activeWorkerCount.Load()))

	}
}
