package task

type task struct {
	arg     interface{}
	handler func(interface{})
	key     int64
}

func (this *task) reset() {
	this.key = 0
	this.arg = nil
	this.handler = nil
}
