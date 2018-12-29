package task4go

import (
	"github.com/smartwalle/pool4go"
	"math"
	"sync"
)

type TaskPool struct {
	maxWorker int
	isRunning bool
	mu        sync.Mutex

	workerPool *pool4go.Pool

	taskList chan func()

	done chan struct{}
}

func NewTaskPool(maxWorker int) *TaskPool {
	var p = &TaskPool{}
	p.maxWorker = maxWorker

	p.taskList = make(chan func(), math.MaxInt32)

	p.run()

	return p
}

func (this *TaskPool) addWorker(w *worker) {
	if this.workerPool != nil {
		this.workerPool.Release(w, false)
	}
}

func (this *TaskPool) getWorker() *worker {
	if this.workerPool == nil {
		return nil
	}
	var conn, err = this.workerPool.Get()
	if err != nil || conn == nil {
		return nil
	}
	return conn.(*worker)
}

func (this *TaskPool) AddTask(task func()) {
	if task == nil {
		return
	}

	select {
	case this.taskList <- task:
	default:
	}
}

func (this *TaskPool) Run() {
	this.run()
}

func (this *TaskPool) run() {
	this.mu.Lock()
	if this.isRunning {
		this.mu.Unlock()
		return
	}

	this.isRunning = true
	this.workerPool = pool4go.NewPool(func() (pool4go.Conn, error) {
		var w = newWorker(this)
		w.start()
		return w, nil
	})
	this.workerPool.SetMaxIdleConns(this.maxWorker)
	this.workerPool.SetMaxOpenConns(this.maxWorker)
	this.done = make(chan struct{}, 1)

	this.mu.Unlock()

	go func() {
		for {
			select {
			case t, ok := <-this.taskList:
				if !ok {
					return
				}

				if t != nil {
					var w = this.getWorker()
					if w != nil {
						w.do(t)
					}
				}
			case <-this.done:
				return
			}
		}
	}()
}

func (this *TaskPool) Stop() {
	this.mu.Lock()
	defer this.mu.Unlock()

	if this.isRunning == false {
		return
	}

	close(this.done)
	this.isRunning = false

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
	return len(this.taskList)
}
