package task4go

type ManagerOption func(m *manager)

func WithWorker(worker int) ManagerOption {
	return func(m *manager) {
		m.worker = worker
	}
}

func WithWaiter(waiter Waiter) ManagerOption {
	return func(m *manager) {
		m.waiter = waiter
	}
}

type TaskOption func(task *task)

func WithArg(arg interface{}) TaskOption {
	return func(task *task) {
		task.arg = arg
	}
}
