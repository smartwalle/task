package main

import (
	"runtime"
	"fmt"
	"github.com/smartwalle/task4go"
	"time"
)

func main() {
	fmt.Println("程序开始运行，此时只有 1 个 Goroutine：", runtime.NumGoroutine())
	var p = task4go.NewTaskPool(5)
	fmt.Println("创建了 TaskPool，此时至少有 2 个 Goroutine：", runtime.NumGoroutine())
	fmt.Println("往 TaskPool 里面添加任务中...")
	for i:=0; i<10000000; i++ {
		p.AddTask(Do)
	}
	fmt.Println("添加了 10000000 个任务")
	fmt.Println("由于 TaskPool 设置的 max worker 为 5，所以此时 Goroutine 的数量为 7：", runtime.NumGoroutine())
	p.Stop()
	time.Sleep(time.Second*5)
	fmt.Println("TaskPool 停止，主 Goroutine 暂停了 5 秒，理论上 TaskPool 的 worker 都已结束运行，所以 Goroutine 的数量为 1：", runtime.NumGoroutine())

	p.Run()
	fmt.Println("重新启动 TaskPool，此时至少有 2 个 Goroutine：", runtime.NumGoroutine())
	time.Sleep(time.Second*5)
	fmt.Println("主 Goroutine 暂停了 5 秒，理论上 TaskPool 的 5 个 worker 都已开始运行，所以 Goroutine 的数量为 7：", runtime.NumGoroutine())
	time.Sleep(time.Second*5)
	p.Stop()
	time.Sleep(time.Second*5)
	fmt.Println("TaskPool 停止，主 Goroutine 暂停了 5 秒，理论上 TaskPool 的 worker 都已结束运行，所以 Goroutine 的数量为 1：", runtime.NumGoroutine())

	p.Run()
	fmt.Println("再次启动 TaskPool，此时至少有 2 个 Goroutine：", runtime.NumGoroutine())
	time.Sleep(time.Second*5)
	fmt.Println("主 Goroutine 暂停了 5 秒，理论上 TaskPool 的 5 个 worker 都已开始运行，所以 Goroutine 的数量为 7：", runtime.NumGoroutine())

	p.SetMaxWorker(3)
	time.Sleep(time.Second * 5)
	fmt.Println("将 TaskPool 的 max worker 降到 3 个，所以 Goroutine 的数量为 5：", runtime.NumGoroutine())

	p.SetMaxWorker(6)
	time.Sleep(time.Second * 5)
	fmt.Println("将 TaskPool 的 max worker 升到 6 个，所以 Goroutine 的数量为 8：", runtime.NumGoroutine())
}

func Do() {
	for i := 0; i < 1000000; i++ {
	}
}