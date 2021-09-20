package task4go

import "sync"

type worker struct {
	index int
	task  chan *task
	pool  *sync.Pool
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

			if t != nil {
				if t.fn != nil {
					t.fn(t.arg)
				}

				t.reset()
				this.pool.Put(t)
			}
		}
	}
}
