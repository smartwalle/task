package main

import (
	"fmt"
	"github.com/smartwalle/task"
	"sync"
	"time"
)

func main() {
	var waiter = &sync.WaitGroup{}
	var m = task.New(task.WithWaiter(waiter), task.WithWorker(10))
	m.Run()

	go func() {
		var i = 0
		for {
			i++
			m.AddTask(func(arg interface{}) {
				fmt.Println("hello", arg)
			}, task.WithArg(i))

			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.AfterFunc(time.Second*10, func() {
		m.Close()
		m.Close()
		m.Close()
		m.Close()
		m.Close()
	})

	fmt.Println("wait..")
	waiter.Wait()
	fmt.Println("done..")
}
