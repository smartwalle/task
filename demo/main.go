package main

import (
	"github.com/smartwalle/task4go"
	"fmt"
	"time"
	"sync"
)

func main() {
	var w = &sync.WaitGroup{}
	var t = time.Now()
	var p = task4go.NewTaskPool(5)
	p.Run()
	p.Run()
	p.Run()
	for i:=0;i <1000000; i++ {
		w.Add(1)
		p.AddTask(&Job{i: i, w: w})
	}
	w.Wait()
	fmt.Println(time.Now().Sub(t))
}

type Job struct {
	i int
	w *sync.WaitGroup
}

func (this *Job) Do() {
	//for i := 0; i < 1000000; i++ {
	//	//fmt.Println("ee")
	//}
	this.w.Done()
	//fmt.Println("do job", this.i)
	//time.Sleep( time.Second * 3)
}