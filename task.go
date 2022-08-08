package task

type task struct {
	key     int64
	arg     interface{}
	handler func(interface{})
}

func (this *task) reset() {
	this.key = 0
	this.arg = nil
	this.handler = nil
}
