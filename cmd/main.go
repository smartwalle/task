package main

import (
	"fmt"
	"github.com/smartwalle/task4go"
	"sync"
	"time"
)

func main() {
	var waiter = &sync.WaitGroup{}
	var p = task4go.New(task4go.WithWaiter(waiter))
	p.Run()

	go func() {
		var i = 0
		for {
			i++
			p.Add(func(arg interface{}) {
				fmt.Println("hello", arg)
			}, task4go.WithArg(i))

			time.Sleep(time.Millisecond * 100)
		}
	}()

	time.AfterFunc(time.Second*10, func() {
		p.Close()
		p.Close()
		p.Close()
		p.Close()
		p.Close()
	})

	fmt.Println("wait..")
	waiter.Wait()
	fmt.Println("done..")
}
