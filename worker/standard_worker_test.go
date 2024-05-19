package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/domain"
	"github.com/vd09/gr_worker/logger"
	"github.com/vd09/gr_worker/worker"
)

const mockTaskDefaultValue = 7

type mockTask struct {
	executed bool
	count    int
}

func (mt *mockTask) ExecuteTask() error {
	mt.executed = true
	mt.count++
	return nil
}

func TestStandardWorker_Start(t *testing.T) {
	// Mock dependencies
	ctx, cancelCtx := context.WithCancel(context.Background())
	mockTask := &mockTask{}
	mockTasks := gr_variable.NewGrChannel[*gr_worker.Task]()

	// Create a StandardWorker instance
	newWorker := worker.NewStandardWorker(ctx, mockTasks, logger.Discard, func(domain.WorkerStatus) bool { return true })

	// Run the worker
	go newWorker.Start()

	// Write a mock task to the tasks channel
	mockTasks.MustWriteValue(gr_worker.NewTask(mockTask.ExecuteTask))
	mockTasks.MustWriteValue(gr_worker.NewTask(mockTask.ExecuteTask))

	// Ensure task is executed
	<-time.After(time.Millisecond) // Wait for some time for execution
	if !mockTask.executed {
		t.Error("Task is not executed")
	}
	if mockTask.count != 2 {
		t.Error("Its not executed 2 times")
	}

	// Ensure worker stops when context is canceled
	cancelCtx()
	mockTasks.StopWriting()

	// Ensure worker stops when tasks channel is closed
	if _, ok := <-mockTasks.Receive(); ok {
		t.Error("Tasks channel not closed")
	}
}
