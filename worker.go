package task

import (
	"sync"
)

type worker struct {
	task  chan *task
	pool  *sync.Pool
	index int
}

func newWorker(index int, task chan *task, pool *sync.Pool) *worker {
	var w = &worker{}
	w.index = index
	w.task = task
	w.pool = pool
	return w
}

func (this *worker) run() {
RunLoop:
	for {
		select {
		case t, ok := <-this.task:
			if ok == false {
				break RunLoop
			}

			if t == nil {
				continue
			}

			if t.handler != nil {
				t.handler(t.arg)
			}

			t.reset()
			this.pool.Put(t)
		}
	}
}
