package task4go

type taskPool struct {
	workerList chan *worker
	taskList   chan Task
}

func NewTaskPool() *taskPool {
	var p = &taskPool{}
	p.workerList = make(chan *worker, 10)
	p.taskList = make(chan Task)

	for i := 0; i < 10; i++ {
		var w = NewWorker()
		p.addWorker(w)
		w.Start()
	}

	p.Run()


	return p
}

func (this *taskPool) addWorker(w *worker) {
	w.pool = this
	this.workerList <- w
}

func (this *taskPool) getWorker() *worker {
	var w = <-this.workerList
	return w
}

func (this *taskPool) AddTask(task Task) {
	this.taskList <- task
}

func (this *taskPool) Run() {
	go func() {
		for {
			select {
			case task := <-this.taskList:
				var w = this.getWorker()
				w.Do(task)
			}
		}
	}()
}
