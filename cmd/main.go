package main

import (
	"fmt"
	"github.com/smartwalle/task4go"
	"sync"
	"time"
)

func main() {
	var waiter = &sync.WaitGroup{}
	var m = task4go.New(task4go.WithWaiter(waiter))
	m.Run()

	go func() {
		var i = 0
		for {
			i++
			m.AddTask(func(arg interface{}) {
				fmt.Println("hello", arg)
			}, task4go.WithArg(i))

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
