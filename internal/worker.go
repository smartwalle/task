package internal

type Worker struct {
	index int
	task  chan *Task
}

func NewWorker(index int, task chan *Task) *Worker {
	var w = &Worker{}
	w.index = index
	w.task = task
	return w
}

func (this *Worker) Run() {
RunLoop:
	for {
		select {
		case task, ok := <-this.task:
			if ok == false {
				break RunLoop
			}

			if task != nil && task.fn != nil {
				task.fn(task.payload)
			}
		}
	}
}
