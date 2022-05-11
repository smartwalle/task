package task4go

type Option func(m *option)

type option struct {
	worker int
	waiter Waiter
}

func WithWorker(worker int) Option {
	return func(m *option) {
		m.worker = worker
	}
}

func WithWaiter(waiter Waiter) Option {
	return func(m *option) {
		m.waiter = waiter
	}
}

type TaskOption func(task *task)

func WithArg(arg interface{}) TaskOption {
	return func(task *task) {
		task.arg = arg
	}
}

func WithKey(key int64) func(task *task) {
	return func(task *task) {
		task.key = key
	}
}
