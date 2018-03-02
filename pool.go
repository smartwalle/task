package task4go

import (
	"github.com/smartwalle/container/slist"
	"runtime"
	"sync/atomic"
)

type taskPool struct {
	maxWorker    int32
	workerNumber int32
	workerList   slist.List
	workerChan   chan struct{}
	taskList     chan Task
}

func NewTaskPool() *taskPool {
	var p = &taskPool{}
	p.maxWorker = int32(runtime.NumCPU())
	p.workerList = slist.New()
	p.workerChan = make(chan struct{})
	p.taskList = make(chan Task)

	p.Run()

	return p
}

func (this *taskPool) addWorker(w *worker) {
	this.workerList.PushBack(w)

	var workerNumber = atomic.LoadInt32(&this.workerNumber)
	var maxWorker = atomic.LoadInt32(&this.maxWorker)
	if workerNumber >= maxWorker {
		this.workerChan <- struct{}{}
	}
}

func (this *taskPool) getWorker() *worker {
	var w = this.workerList.PopFront()
	if w != nil {
		return w.(*worker)
	}

	var workerNumber = atomic.LoadInt32(&this.workerNumber)
	var maxWorker = atomic.LoadInt32(&this.maxWorker)

	if workerNumber >= maxWorker {
		<-this.workerChan
		return this.getWorker()
	}

	var nw = NewWorker(this)
	nw.Start()
	atomic.StoreInt32(&this.workerNumber, this.workerNumber)

	return nw
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
