package task4go

import (
	"github.com/smartwalle/container/slist"
	"github.com/smartwalle/pool4go"
	"math"
	"sync"
)

type TaskPool struct {
	maxWorker int
	running   bool
	mu        sync.Mutex

	workerPool *pool4go.Pool

	taskEvent chan struct{}
	taskList  slist.List

	stopEvent chan struct{}
}

func NewTaskPool(maxWorker int) *TaskPool {
	var p = &TaskPool{}
	p.maxWorker = maxWorker

	p.taskList = slist.New()
	p.taskEvent = make(chan struct{}, math.MaxInt32)

	p.run()

	return p
}

func (this *TaskPool) addWorker(w *worker) {
	this.workerPool.Release(w, false)
}

func (this *TaskPool) getWorker() *worker {
	var conn, err = this.workerPool.Get()
	if err != nil {
		return nil
	}
	return conn.(*worker)
}

func (this *TaskPool) AddTask(task func()) {
	if task == nil {
		return
	}
	this.taskList.PushBack(task)
	this.taskEvent <- struct{}{}
}

func (this *TaskPool) Run() {
	this.run()
}

func (this *TaskPool) run() {
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
	})
	this.workerPool.SetMaxIdleConns(this.maxWorker)
	this.workerPool.SetMaxOpenConns(this.maxWorker)
	this.stopEvent = make(chan struct{})

	this.mu.Unlock()

	go func() {
		for {
			select {
			case <-this.taskEvent:
				var t = this.taskList.PopFront()
				if t != nil {
					var w = this.getWorker()
					if w != nil {
						w.do(t.(func()))
					}
				}
			case <-this.stopEvent:
				return
			}
		}
	}()
}

func (this *TaskPool) Stop() {
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
	this.workerPool = nil
}

func (this *TaskPool) SetMaxWorker(n int) {
	this.maxWorker = n
	if this.workerPool != nil {
		this.workerPool.SetMaxOpenConns(this.maxWorker)
		this.workerPool.SetMaxIdleConns(this.maxWorker)
	}
}

func (this *TaskPool) MaxWorker() int {
	return this.maxWorker
}

func (this *TaskPool) NumTask() int {
	return this.taskList.Len()
}
