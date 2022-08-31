package task

import (
	"errors"
	"github.com/smartwalle/queue/block"
	"sync"
	"sync/atomic"
)

var (
	ErrClosed  = errors.New("task manager closed")
	ErrBadTask = errors.New("bad task")
)

type Manager interface {
	Run()

	Close()

	AddTask(handler func(arg interface{}), opts ...TaskOption) error
}

type manager struct {
	options  *options
	pool     *sync.Pool
	queue    block.Queue[*task]
	dispatch chan *task
	runOnce  sync.Once
	closed   int32
}

func New(opts ...Option) Manager {
	var pool = &sync.Pool{
		New: func() interface{} {
			return &task{}
		},
	}
	return newManager(pool, opts...)
}

func newManager(pool *sync.Pool, opts ...Option) *manager {
	var m = &manager{}
	m.options = &options{}
	m.pool = pool
	m.queue = block.New[*task]()
	m.dispatch = make(chan *task, 1)
	m.closed = 0

	for _, opt := range opts {
		if opt != nil {
			opt(m.options)
		}
	}

	if m.options.worker <= 0 {
		m.options.worker = 1
	}
	if m.options.waiter == nil {
		m.options.waiter = &sync.WaitGroup{}
	}
	return m
}

func (this *manager) Run() {
	this.runOnce.Do(this.run)
}

func (this *manager) run() {
	for i := 1; i <= this.options.worker; i++ {
		this.options.waiter.Add(1)
		var w = newWorker(i, this.dispatch, this.pool)
		go func() {
			w.run()
			this.options.waiter.Done()
		}()
	}

	go func() {
		var nTasks []*task
	RunLoop:
		for {
			nTasks = nTasks[0:0]
			var ok = this.queue.Dequeue(&nTasks)

			for _, nTask := range nTasks {
				this.dispatch <- nTask
			}

			if ok == false {
				break RunLoop
			}
		}
		close(this.dispatch)
	}()
}

func (this *manager) Close() {
	if atomic.CompareAndSwapInt32(&this.closed, 0, 1) {
		this.queue.Close()
	}
}

func (this *manager) AddTask(handler func(arg interface{}), opts ...TaskOption) error {
	if handler == nil {
		return ErrBadTask
	}

	if atomic.LoadInt32(&this.closed) == 1 {
		return ErrClosed
	}

	var nTask = this.pool.Get().(*task)
	nTask.handler = handler

	for _, opt := range opts {
		if opt != nil {
			opt(nTask)
		}
	}

	if this.queue.Enqueue(nTask) == false {
		return ErrClosed
	}
	return nil
}
