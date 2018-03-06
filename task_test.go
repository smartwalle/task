package task4go

import (
	"testing"
)

func Benchmark_Do(b *testing.B) {
	var p = NewTaskPool(5)
	for i:=0; i<b.N; i++ {
		p.AddTask(Do)
	}
}

func Do() {
	for i := 0; i < 1000000; i++ {
	}
}