package main

import (
	"fmt"
	"github.com/smartwalle/task4go"
	"sync"
	"time"
)

func main() {
	var waiter = &sync.WaitGroup{}
	var p = task4go.New(2, waiter)
	p.Run()

	go func() {
		for {
			p.Add(func(payload interface{}) {
				fmt.Println("hahaha")
			}, nil)

			time.Sleep(time.Second)
		}
	}()

	time.AfterFunc(time.Second*10, func() {
		p.Close()
	})

	fmt.Println("wait..")
	waiter.Wait()
	fmt.Println("done..")
}
