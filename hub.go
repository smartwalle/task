package task

import (
	"sync"
	"sync/atomic"
)

type hub struct {
	pool    *sync.Pool
	manages []*manager
	worker  int64
	runOnce sync.Once
}

func NewHub(opts ...Option) Manager {
	var nOpts = &options{}

	for _, opt := range opts {
		if opt != nil {
			opt(nOpts)
		}
	}

	if nOpts.worker <= 0 {
		nOpts.worker = 1
	}
	if nOpts.waiter == nil {
		nOpts.waiter = &sync.WaitGroup{}
	}

	var pool = &sync.Pool{
		New: func() interface{} {
			return &task{}
		},
	}

	var nHub = &hub{}
	nHub.worker = int64(nOpts.worker)
	nHub.manages = make([]*manager, nHub.worker)
	nHub.pool = pool
	for idx := range nHub.manages {
		nHub.manages[idx] = newManager(nHub.pool, WithWorker(1), WithWaiter(nOpts.waiter))
	}

	return nHub
}

func (this *hub) Run() {
	this.runOnce.Do(func() {
		for _, m := range this.manages {
			m.run()
		}
	})
}

func (this *hub) Close() {
	for _, m := range this.manages {
		m.Close()
	}
}

func (this *hub) AddTask(handler func(arg interface{}), opts ...TaskOption) error {
	if handler == nil {
		return ErrBadTask
	}

	var nTask = this.pool.Get().(*task)
	nTask.handler = handler

	for _, opt := range opts {
		if opt != nil {
			opt(nTask)
		}
	}

	var idx = nTask.key % this.worker
	if idx < 0 {
		idx *= -1
	}
	var m = this.manages[idx]

	if atomic.LoadInt32(&m.closed) == 1 {
		return ErrClosed
	}

	if m.queue.Enqueue(nTask) == false {
		return ErrClosed
	}

	return nil
}
