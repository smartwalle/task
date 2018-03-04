package task4go

import (
	"github.com/smartwalle/container/slist"
	"math"
	"sync"
)

type taskPool struct {
	maxWorker int
	running   bool
	mu sync.Mutex

	workerList chan *worker

	taskEvent chan struct{}
	taskList  slist.List

	stopEvent       chan struct{}
	stopWorkerEvent chan chan struct{}
}

func NewTaskPool(maxWorker int) *taskPool {
	var p = &taskPool{}
	p.maxWorker = maxWorker

	p.taskList = slist.New()
	p.taskEvent = make(chan struct{}, math.MaxInt32)

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

func (this *taskPool) AddTask(task func()) {
	if task == nil {
		return
	}
	this.taskList.PushBack(task)
	this.taskEvent <- struct{}{}
}

func (this *taskPool) Run() {
	this.run()
}

func (this *taskPool) run() {
	this.mu.Lock()
	if this.running {
		this.mu.Unlock()
		return
	}

	this.running = true
	this.workerList = make(chan *worker, this.maxWorker)
	this.stopEvent = make(chan struct{})
	this.stopWorkerEvent = make(chan chan struct{}, this.maxWorker)

	for i := 0; i < this.maxWorker; i++ {
		var w = newWorker(this)
		w.start()
		this.addWorker(w)
	}
	this.mu.Unlock()

	go func() {
		for {
			select {
			case <-this.taskEvent:
				var w = this.getWorker()
				var t = this.taskList.PopFront()
				if t != nil {
					w.do(t.(func()))
				}
			case <-this.stopEvent:
				this.running = false
				return
			}
		}
	}()
}

func (this *taskPool) Stop() {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.running == false {
		return
	}

	this.stopEvent <- struct{}{}

	for i := 0; i < this.maxWorker; i++ {
		var w = this.getWorker()
		w.stop()
	}
	close(this.workerList)
	close(this.stopWorkerEvent)
}
