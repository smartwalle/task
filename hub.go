package task

import (
	"sync"
	"sync/atomic"
)

type hub struct {
	worker  int64
	manages []*manager
	pool    *sync.Pool
	runOnce sync.Once
}

func NewHub(opts ...Option) Manager {
	var nOpt = &option{}

	for _, opt := range opts {
		if opt != nil {
			opt(nOpt)
		}
	}

	if nOpt.worker <= 0 {
		nOpt.worker = 1
	}
	if nOpt.waiter == nil {
		nOpt.waiter = &sync.WaitGroup{}
	}

	var pool = &sync.Pool{
		New: func() interface{} {
			return &task{}
		},
	}

	var nHub = &hub{}
	nHub.worker = int64(nOpt.worker)
	nHub.manages = make([]*manager, nHub.worker)
	nHub.pool = pool
	for idx := range nHub.manages {
		nHub.manages[idx] = newManager(nHub.pool, WithWorker(1), WithWaiter(nOpt.waiter))
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

func (this *hub) AddTask(fn func(arg interface{}), opts ...TaskOption) error {
	if fn == nil {
		return ErrBadTask
	}

	var nTask = this.pool.Get().(*task)
	nTask.fn = fn

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
