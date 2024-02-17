package entities

type Actor interface {
	AddTask(task Task) error
	Start()
	Stop()
}
