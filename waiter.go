package task4go

type Waiter interface {
	Add(delta int)

	Done()

	Wait()
}
