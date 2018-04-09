package task4go

import (
	"os"
	"testing"
)

var p *TaskPool

func TestMain(m *testing.M) {
	p = NewTaskPool(5)

	os.Exit(m.Run())
}

func Benchmark_Do(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p.AddTask(Do)
	}
}

func Do() {
	for i := 0; i < 1000000; i++ {
	}
}
