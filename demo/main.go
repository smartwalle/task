package main

import (
	"github.com/smartwalle/task4go"
	"fmt"
	"time"
)

func main() {
	var p = task4go.NewTaskPool(5)
	for i:=0;i <1000000; i++ {
		p.AddTask(&Job{})
	}
	fmt.Println("end")
	time.Sleep(time.Second * 10)
}

type Job struct {
}

func (this *Job) Do() {
	fmt.Println("do job")
	//time.Sleep( time.Second * 3)
}
