package task4go

type taskPool struct {
	maxWorker    int
	workerList   chan *worker
	taskList     chan Task
}

func NewTaskPool(maxWorker int) *taskPool {
	var p = &taskPool{}
	p.maxWorker = maxWorker
	p.workerList = make(chan *worker, maxWorker)
	p.taskList = make(chan Task)

	for i:=0; i<maxWorker; i++ {
		var w = NewWorker(p)
		p.addWorker(w)
	}
	p.run()

	return p
}

func (this *taskPool) addWorker(w *worker) {
	this.workerList <- w
}

func (this *taskPool) getWorker() *worker {
	var w = <- this.workerList
	return w
}

func (this *taskPool) AddTask(task Task) {
	this.taskList <- task
}

func (this *taskPool) run() {
	go func() {
		for {
			select {
			case task := <-this.taskList:
				go func(task Task) {
					var w = this.getWorker()
					w.Do(task)
				}(task)
			}
		}
	}()
}
