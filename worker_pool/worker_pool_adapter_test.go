package worker_pool

import (
	"fmt"
	"math"
	"sync/atomic"
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

	for i := 0; i < 5; i++ {
		taskAdded := wp.AddTask(basicFunction)
		time.Sleep(time.Millisecond)
		if !taskAdded {
			t.Error("task not added to worker pool")
		}
		workers := math.Min(float64(i+1), 3)
		if wp.activeWorkerCount.Load() != int32(workers) {
			t.Error(fmt.Sprintf("total worker (%d) are not matching as expected (%v)", wp.activeWorkerCount.Load(), workers))
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

func TestWorkerPoolAdapter_WaitAndStop(t *testing.T) {
	// Create a worker pool with default options
	wp, err := NewWorkerPool(
		WithMinWorkers(1),
		WithMaxWorkers(3),
		WithMaxTasks(5),
	)
	if err != nil {
		t.Fatalf("error creating worker pool: %v", err)
	}

	// Define a task that will take some time to complete
	taskDuration := 2 * time.Second
	var count atomic.Int32
	basicFunction := func() {
		time.Sleep(taskDuration)
		count.Add(1)
	}

	// Add tasks to the worker pool
	numTasks := 3
	for i := 0; i < numTasks; i++ {
		taskAdded := wp.AddTask(basicFunction)
		if !taskAdded {
			t.Errorf("task %d not added to worker pool", i+1)
		}
	}

	// Call WaitAndStop and measure the time taken
	wp.WaitAndStop()

	if int(count.Load()) != numTasks {
		t.Errorf("WaitAndStop returned before all tasks completed; elapsed tasks: %v, expected: %v", count.Load(), numTasks)
	}

	// Verify that the worker pool is stopped
	if !wp.IsWorkerPoolStopped() {
		t.Error("worker pool not stopped after WaitAndStop")
	}
}
