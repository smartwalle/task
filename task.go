package task4go

import (
	"github.com/smartwalle/container/slist"
	"github.com/smartwalle/pool4go"
	"math"
	"sync"
)

type taskPool struct {
	maxWorker int
	running   bool
	mu        sync.Mutex

	workerPool *pool4go.Pool

	taskEvent chan struct{}
	taskList  slist.List

	stopEvent chan struct{}
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
	this.workerPool.Release(w, false)
}

func (this *taskPool) getWorker() *worker {
	var conn, err = this.workerPool.Get()
	if err != nil {
		return nil
	}
	return conn.(*worker)
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
	this.workerPool = pool4go.NewPool(func() (pool4go.Conn, error) {
		var w = newWorker(this)
		w.start()
		return w, nil
	}, this.maxWorker)
	this.stopEvent = make(chan struct{})

	this.mu.Unlock()

	go func() {
		for {
			select {
			case <-this.taskEvent:
				var w = this.getWorker()
				var t = this.taskList.PopFront()
				if t != nil && w != nil {
					w.do(t.(func()))
				}
			case <-this.stopEvent:
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
	this.running = false

	for i := 0; i < this.maxWorker; i++ {
		var w = this.getWorker()
		w.stop()
	}
	this.workerPool.Close()
}
