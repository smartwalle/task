package task4go

type task struct {
	arg interface{}
	fn  func(interface{})
}

func (this *task) reset() {
	this.arg = nil
	this.fn = nil
}
