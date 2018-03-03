package main

import (
	"testing"
	"github.com/smartwalle/task4go"
)

func BenchmarkJob_Do(b *testing.B) {
	var p = task4go.NewTaskPool(5)
	for i:=0; i<b.N; i++ {
		p.AddTask(&Job{i: i})
	}
}