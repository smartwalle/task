package internal

import (
	"sync"
)

type Queue struct {
	cond  *sync.Cond
	tasks []*Task
}

func NewQueue() *Queue {
	var q = &Queue{}
	q.cond = sync.NewCond(&sync.Mutex{})
	q.tasks = make([]*Task, 0, 128)
	return q
}

func (this *Queue) Enqueue(task *Task) {
	this.cond.L.Lock()
	this.tasks = append(this.tasks, task)
	this.cond.L.Unlock()

	this.cond.Signal()
}

func (this *Queue) Dequeue(tasks *[]*Task) {
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

func (this *Queue) Close() {

}
