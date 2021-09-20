package task4go

import (
	"sync"
	"sync/atomic"
)

type Manager interface {
	Run()

	Close()

	Add(fn func(arg interface{}), opts ...TaskOption)
}

type manager struct {
	worker   int
	waiter   Waiter
	queue    *queue
	pool     *sync.Pool
	dispatch chan *task
	runOnce  sync.Once
	closed   int32
}

func New(opts ...ManagerOption) Manager {
	var m = &manager{}
	m.queue = newQueue()
	m.pool = &sync.Pool{
		New: func() interface{} {
			return &task{}
		},
	}
	m.dispatch = make(chan *task, 1)
	m.closed = 0

	for _, opt := range opts {
		opt(m)
	}

	if m.worker <= 0 {
		m.worker = 1
	}
	if m.waiter == nil {
		m.waiter = &sync.WaitGroup{}
	}
	return m
}

func (this *manager) Run() {
	this.runOnce.Do(this.run)
}

func (this *manager) run() {
	for i := 0; i < this.worker; i++ {
		this.waiter.Add(1)
		var w = newWorker(i+1, this.dispatch, this.pool)
		go func() {
			w.run()
			this.waiter.Done()
		}()
	}

	go func() {
		var nTasks []*task
	RunLoop:
		for {
			nTasks = nTasks[0:0]
			this.queue.dequeue(&nTasks)

			for _, nTask := range nTasks {
				if nTask == nil {
					break RunLoop
				}
				this.dispatch <- nTask
			}
		}
		close(this.dispatch)
	}()
}

func (this *manager) Close() {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		this.queue.enqueue(nil)
	}
}

func (this *manager) Add(fn func(arg interface{}), opts ...TaskOption) {
	if fn == nil {
		return
	}

	if atomic.LoadInt32(&this.closed) == 1 {
		return
	}

	var nTask, _ = this.pool.Get().(*task)
	nTask.fn = fn
	nTask.arg = nil

	for _, opt := range opts {
		opt(nTask)
	}

	this.queue.enqueue(nTask)
}
