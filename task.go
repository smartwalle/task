package task4go

type task struct {
	key int64
	arg interface{}
	fn  func(interface{})
}

func (this *task) reset() {
	this.key = 0
	this.arg = nil
	this.fn = nil
}
