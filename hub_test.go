package task4go_test

import (
	"github.com/smartwalle/task4go"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestHub_Close(t *testing.T) {
	var waiter = &sync.WaitGroup{}
	var m = task4go.NewHub(task4go.WithWorker(3), task4go.WithWaiter(waiter))
	m.Run()

	m.AddTask(func(arg interface{}) {
		t.Logf("[%d]Task %v start...", time.Now().Unix(), arg)
		time.Sleep(time.Second * 2)
		t.Logf("[%d]Task %v done.", time.Now().Unix(), arg)
	}, task4go.WithArg(1), task4go.WithKey(1))

	m.AddTask(func(arg interface{}) {
		t.Logf("[%d]Task %v start...", time.Now().Unix(), arg)
		time.Sleep(time.Second * 5)
		t.Logf("[%d]Task %v done.", time.Now().Unix(), arg)
	}, task4go.WithArg(2), task4go.WithKey(2))

	go func() {
		time.Sleep(time.Second * 2)
		t.Logf("[%d]Close.", time.Now().Unix())
		m.Close()
	}()

	waiter.Wait()
	t.Logf("[%d]Done.", time.Now().Unix())
}

func TestHub_Run(t *testing.T) {
	var waiter = &sync.WaitGroup{}
	var m = task4go.NewHub(task4go.WithWorker(3), task4go.WithWaiter(waiter))
	m.Run()

	var s = 8000000
	var c = int32(0)
	for i := 0; i < s; i++ {
		m.AddTask(func(arg interface{}) {
			atomic.AddInt32(&c, 1)
		}, task4go.WithKey(int64(i)))
	}

	m.Close()
	waiter.Wait()

	if int(c) != s {
		t.Fatal("任务未执行完成")
	}
}

func TestHub_Run2(t *testing.T) {
	var worker = 3
	var waiter = &sync.WaitGroup{}
	var m = task4go.NewHub(task4go.WithWorker(worker), task4go.WithWaiter(waiter))
	m.Run()

	var r1 = make([]int, worker)
	var r2 = make([]int, worker)

	for i := 0; i < 10000; i++ {
		r1[i%worker] += i

		m.AddTask(func(arg interface{}) {
			var v = arg.(int)
			var idx = v % worker
			r2[idx] += v
		}, task4go.WithArg(i), task4go.WithKey(int64(i)))
	}

	m.Close()
	waiter.Wait()

	for i := 0; i < worker; i++ {
		if r1[i] != r2[i] {
			t.Log("计算结果不匹配")
		}
	}
}
