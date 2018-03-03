package main

import (
	"github.com/smartwalle/task4go"
	"fmt"
	"time"
)

func main() {
	var t = time.Now()
	var p = task4go.NewTaskPool(5)
	for i:=0;i <1000000; i++ {
		p.AddTask(&Job{i: i})
	}

	fmt.Println(time.Now().Sub(t))
}

type Job struct {
	i int
}

func (this *Job) Do() {
	for i := 0; i < 1000000; i++ {
	}
	//this.w.Done()
	//fmt.Println("do job", this.i)
	//time.Sleep( time.Second * 3)
}
