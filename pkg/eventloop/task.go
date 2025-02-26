package eventloop

type Task struct {
	MainTask   func()
	Callback   func()
	IsBlocking bool
}
