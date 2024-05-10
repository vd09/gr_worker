package worker

import (
	"context"
	"fmt"
	"testing"
	"time"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
	"github.com/vd09/gr_worker/logger"
)

const mockTaskDefaultValue = 7

type mockTask struct {
	executed bool
	data     []int
}

func (mt *mockTask) ExecuteTask() error {
	mt.executed = true
	mt.data = append(mt.data, mockTaskDefaultValue)
	return nil
}

func TestSingleTaskWorker_Start(t *testing.T) {
	// Mock dependencies
	ctx, cancelCtx := context.WithCancel(context.Background())
	mockTask := &mockTask{}
	mockTasks := gr_variable.NewGrChannel[*gr_worker.Task]()

	// Create a SingleTaskWorker instance
	worker := NewSingleTaskWorker(ctx, mockTasks, logger.Discard)

	// Run the worker
	go worker.Start()
	mockTasks.MustWriteValue(gr_worker.NewTask(mockTask.ExecuteTask))

	// Ensure task is executed
	<-time.After(time.Millisecond) // Wait for some time for execution
	if !mockTask.executed {
		t.Error("Task is not executed")
	}

	// Ensure worker stops when context is canceled
	cancelCtx()
	mockTasks.StopWriting()
	select {
	case _, ok := <-mockTasks.Receive():
		if ok {
			t.Error("Tasks channel not closed")
		}
	default:
		t.Error("Tasks channel not closed")
	}

	if len(mockTask.data) <= 1 {
		t.Error(fmt.Sprintf("data is not present while executing; data:%#v", mockTask.data))
	}

	for _, val := range mockTask.data {
		if val != mockTaskDefaultValue {
			t.Error(fmt.Sprintf("data (%#v) is not matching with expected data (%#v)", mockTask.data, mockTaskDefaultValue))
		}
	}
}
