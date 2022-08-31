package task

type Option func(m *options)

type options struct {
	waiter Waiter
	worker int
}

func WithWorker(worker int) Option {
	return func(opts *options) {
		opts.worker = worker
	}
}

func WithWaiter(waiter Waiter) Option {
	return func(opts *options) {
		opts.waiter = waiter
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
