package task4go

type worker struct {
	pool *taskPool
	task chan func()
	done chan struct{}
}

func newWorker(pool *taskPool) *worker {
	var w = &worker{}
	w.pool = pool
	w.task = make(chan func())
	w.done = make(chan struct{})

	//pool.stopWorkerEvent <- w.done

	return w
}

func (this *worker) start() {
	go func() {
		for {
			select {
			case t := <-this.task:
				if t != nil {
					t()
				}

				if this.pool != nil {
					this.pool.addWorker(this)
				}
			case <-this.done:
				return
			}
		}
	}()
}

func (this *worker) stop() {
	this.done <- struct{}{}
}

func (this *worker) do(task func()) {
	this.task <- task
}
