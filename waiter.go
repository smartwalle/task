package task

type Waiter interface {
	Add(delta int)

	Done()

	Wait()
}
