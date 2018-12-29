package task4go

type worker struct {
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
			case <-this.pool.done:
				return
			}
		}
	}()
}

func (this *worker) do(task func()) {
	select {
	case this.task <- task:
	default:
	}
}

func (this *worker) Close() error {
	close(this.task)
	return nil
}
