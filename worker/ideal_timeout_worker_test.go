package worker

import (
	"context"
	"testing"
	"time"

	gr_variable "github.com/vd09/gr-variable"
	"github.com/vd09/gr_worker"
)

type mockLogger struct{}

func (m *mockLogger) Printf(format string, v ...interface{}) {}

func TestIdealTimeoutWorker_Start(t *testing.T) {
	// Mock dependencies
	ctx := context.Background()
	mockLogger := &mockLogger{}
	mockTasks := gr_variable.NewGrChannel[*gr_worker.Task]()

	isMockStopFuncCalled := false
	mockStopFunc := func(bool) bool { isMockStopFuncCalled = true; return true }

	// Create an IdealTimeoutWorker instance
	worker := NewIdealTimeoutWorker(ctx, time.Second, mockLogger, mockTasks, mockStopFunc)
	go worker.Start()

	// Test when task execution returns an error
	taskWorker := false
	mockTasks.MustWriteValue(gr_worker.NewTask(func() error {
		taskWorker = true
		return nil
	}))

	// Test when timer ticks
	timer := time.NewTimer(2 * time.Second)
	<-timer.C

	if !taskWorker {
		t.Error("Task is not executed")
	}
	if !isMockStopFuncCalled {
		t.Error("can stop worker function is not called")
	}
}
