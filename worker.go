package task4go

import "sync"

type worker struct {
	mu   sync.Mutex
	pool *TaskPool
	task chan func()
}

func newWorker(pool *TaskPool) *worker {
	var w = &worker{}
	w.pool = pool
	w.task = make(chan func(), 1)
	return w
}

func (this *worker) start() {
	go func() {
		for {
			select {
			case t, ok := <-this.task:
				if !ok {
					return
				}

				if t != nil {
					t()
				}

				if this.pool != nil {
					this.pool.addWorker(this)
				}
			case <-this.pool.closeChan:
				return
			}
		}
	}()
}

func (this *worker) do(task func()) {
	this.mu.Lock()
	select {
	case this.task <- task:
	default:
	}
	this.mu.Unlock()
}

func (this *worker) Close() error {
	this.mu.Lock()
	close(this.task)
	this.task = nil
	this.mu.Unlock()
	return nil
}
