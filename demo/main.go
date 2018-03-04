package main

import (
	"github.com/smartwalle/task4go"
	"fmt"
	"time"
	"sync"
)

func main() {
	var w = sync.WaitGroup{}
	var t = time.Now()
	var p = task4go.NewTaskPool(5)
	p.Run()
	p.Run()
	p.Run()
	for i:=0;i <1000000; i++ {
		w.Add(1)
		p.AddTask(&Job{i: i})
	}

	time.Sleep(time.Second * 3)
	p.Stop()

	p.Run()

	for i:=0;i <100; i++ {
		p.AddTask(&Job2{i: i})
	}

	w.Wait()
	fmt.Println(time.Now().Sub(t))
}

type Job struct {
	i int
}

func (this *Job) Do() {
	for i := 0; i < 1000000; i++ {
		//fmt.Println("ee")
	}
	//this.w.Done()
	//fmt.Println("do job", this.i)
	//time.Sleep( time.Second * 3)
}

type Job2 struct {
	i int
}


func (this *Job2) Do() {
	fmt.Println("job2", this.i)
	//this.w.Done()
	//fmt.Println("do job", this.i)
	//time.Sleep( time.Second * 3)
}