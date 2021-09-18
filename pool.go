package task4go

import (
	"github.com/smartwalle/task4go/internal"
	"sync"
	"sync/atomic"
)

type Pool struct {
	worker   int
	queue    *internal.Queue
	dispatch chan *internal.Task
	closed   int32
	waiter   Waiter
}

func New(worker int, waiter Waiter) *Pool {
	if waiter == nil {
		waiter = &sync.WaitGroup{}
	}

	var p = &Pool{}
	p.worker = worker
	p.queue = internal.NewQueue()
	p.dispatch = make(chan *internal.Task, 1)
	p.closed = 0
	p.waiter = waiter
	return p
}

func (this *Pool) Run() {
	for i := 0; i < this.worker; i++ {
		this.waiter.Add(1)
		var w = internal.NewWorker(i+1, this.dispatch)
		go func() {
			w.Run()
			this.waiter.Done()
		}()
	}

	go func() {
		var nTasks []*internal.Task
	RunLoop:
		for {
			nTasks = nTasks[0:0]
			this.queue.Dequeue(&nTasks)

			for _, nTask := range nTasks {
				if nTask == nil {
					break RunLoop
				}
				this.dispatch <- nTask
			}
		}
		atomic.SwapInt32(&this.closed, 1)
		close(this.dispatch)
	}()
}

func (this *Pool) Close() {
	if atomic.LoadInt32(&this.closed) == 1 {
		return
	}
	this.queue.Enqueue(nil)
}

func (this *Pool) Add(fn func(payload interface{}), payload interface{}) {
	if fn == nil {
		return
	}

	if atomic.LoadInt32(&this.closed) == 1 {
		return
	}

	var nTask = internal.NewTask(fn, payload)
	this.queue.Enqueue(nTask)
}
