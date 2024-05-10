package gr_worker

import (
	"errors"
	"fmt"
	"reflect"
)

type Task struct {
	fn     interface{}
	params []interface{}
}

func NewTask(taskFunc interface{}, params ...interface{}) *Task {
	return &Task{
		fn:     taskFunc,
		params: params,
	}
}

// ExecuteTask executes a given task function with parameters.
func (t *Task) ExecuteTask() error {
	task := reflect.ValueOf(t.fn)
	if task.Kind() != reflect.Func {
		return errors.New(fmt.Sprintf("Error: taskFunc %v is not a function", task))
	}

	if len(t.params) < task.Type().NumIn() {
		return errors.New(fmt.Sprintf("Error: number of parameters does not match; expected:%v actual:%v",
			task.Type().NumIn(), len(t.params)))
	}

	inputs := make([]reflect.Value, len(t.params))
	for i, param := range t.params {
		if param == nil {
			inputs[i] = reflect.Zero(task.Type().In(i)) // Pass zero value for nil parameter
		} else {
			inputs[i] = reflect.ValueOf(param)
		}
	}

	// Call the task function and handle the return values
	task.Call(inputs)
	//results := task.Call(inputs)
	//if len(results) > 0 && !results[0].IsNil() {
	//	// If the first return value is not nil, assume it's an error
	//	return results[0].Interface().(error)
	//}
	return nil
}
