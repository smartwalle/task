package internal

type Task struct {
	payload interface{}
	fn      func(interface{})
}

func NewTask(fn func(interface{}), payload interface{}) *Task {
	var t = &Task{}
	t.payload = payload
	t.fn = fn
	return t
}
