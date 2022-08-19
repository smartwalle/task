package task_test

import (
	"github.com/smartwalle/task"
	"sync"
	"testing"
)

func BenchmarkHub_AddTask(b *testing.B) {
	var waiter = &sync.WaitGroup{}
	var m = task.New(task.WithWorker(1), task.WithWaiter(waiter))
	m.Run()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		m.AddTask(func(arg interface{}) {
		})
	}

	m.Close()
	waiter.Wait()
}
