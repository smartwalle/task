package task4go

import (
	"sync"
)

type queue struct {
	cond  *sync.Cond
	tasks []*task
}

func newQueue() *queue {
	var q = &queue{}
	q.cond = sync.NewCond(&sync.Mutex{})
	q.tasks = make([]*task, 0, 128)
	return q
}

func (this *queue) enqueue(t *task) {
	this.cond.L.Lock()
	this.tasks = append(this.tasks, t)
	this.cond.L.Unlock()

	this.cond.Signal()
}

func (this *queue) dequeue(tasks *[]*task) {
	this.cond.L.Lock()
	for len(this.tasks) == 0 {
		this.cond.Wait()
	}

	for _, task := range this.tasks {
		*tasks = append(*tasks, task)
		if task == nil {
			break
		}
	}

	this.tasks = this.tasks[0:0]
	this.cond.L.Unlock()
}
