package task4go

import (
	"github.com/smartwalle/container/slist"
	"math"
)

type taskPool struct {
	maxWorker  int
	workerList chan *worker
	taskEvent  chan struct{}
	taskList   slist.List
	stopEvent  chan struct{}
}

func NewTaskPool(maxWorker int) *taskPool {
	var p = &taskPool{}
	p.maxWorker = maxWorker
	p.workerList = make(chan *worker, maxWorker)
	p.taskEvent = make(chan struct{}, math.MaxInt32)
	p.taskList = slist.New()
	p.stopEvent = make(chan struct{})

	for i := 0; i < maxWorker; i++ {
		var w = NewWorker(p)
		w.Start()
		p.addWorker(w)
	}
	p.run()

	return p
}

func (this *taskPool) addWorker(w *worker) {
	this.workerList <- w
}

func (this *taskPool) getWorker() *worker {
	var w = <-this.workerList
	return w
}

func (this *taskPool) AddTask(task Task) {
	if task == nil {
		return
	}
	this.taskList.PushBack(task)
	this.taskEvent <- struct{}{}
}

func (this *taskPool) run() {
	go func() {
		for {
			select {
			case <-this.taskEvent:
				var w = this.getWorker()
				var t = this.taskList.PopFront()
				if t != nil {
					w.Do(t.(Task))
				}
			case <-this.stopEvent:
				return
			}
		}
	}()
}

func (this *taskPool) Stop() {
	this.stopEvent <- struct{}{}
}
