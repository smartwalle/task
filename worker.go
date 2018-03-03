package task4go

type worker struct {
	pool *taskPool
	task chan Task
	quit chan struct{}
}

func NewWorker(pool *taskPool) *worker {
	var w = &worker{}
	w.pool = pool
	w.task = make(chan Task)
	w.quit = make(chan struct{})
	return w
}

func (this *worker) Start() {
	go func() {
		for {
			select {
			case t := <-this.task:
				if t != nil {
					t.Do()
				}

				if this.pool != nil {
					this.pool.addWorker(this)
				}
			case <-this.quit:
				return
			}
		}
	}()
}

func (this *worker) Stop() {
	this.quit <- struct{}{}
}

func (this *worker) Do(task Task) {
	this.task <- task
}
