package eventloop

import "sync"

type EventLoop struct {
	mainTasks chan Task // channel to hold tasks waiting to be processed
	taskQueue chan Task // channel to hold callback tasks
	stop      chan bool // channel to indicate the event loop to stop
}

func NewEventLoop(taskBufferSize int) *EventLoop {
	return &EventLoop{
		mainTasks: make(chan Task, taskBufferSize), // Initialize the main tasks channel with a buffer size of 10
		taskQueue: make(chan Task, taskBufferSize), // Initialize the tasks queue channel with a buffer size of 10
		stop:      make(chan bool),                 // Initialize the stop channel
	}
}

// Add adds a new task to the mainTasks channel for immediate processing.
func (eventLoop *EventLoop) Add(task *Task) {
	eventLoop.mainTasks <- *task // Push task to mainTasks channel for execution
}

// AddTaskToQueue adds a new blocking task to the taskQueue for later processing.
func (eventLoop *EventLoop) AddTaskToQueue(task *Task) {
	eventLoop.taskQueue <- *task
}

// StopeventLoop sends a stop signal to the event loop, instructing it to terminate.
func (eventLoop *EventLoop) StopEventLoop() {
	eventLoop.stop <- true // send signal to stop the event loop
}

func (eventLoop *EventLoop) InitEventLoop(workerPoolSize int) *sync.WaitGroup {
	var wg sync.WaitGroup

	// add the event loop goroutine to the wait group to trac it's execution
	wg.Add(1)

	//worker pool for the handling blocking tasks
	workerPool := make(chan struct{}, workerPoolSize)

	// start the event loop in a separate goroutine
	go func() {
		defer wg.Done() // Ensure the wait group  is decremented when the goroutine finishes

		// Main event loop : process tasks continuously
		for {
			// wait for tasks from the mainTasks or taskQueue channels
			select {
			case task := <-eventLoop.mainTasks:
				if task.IsBlocking {
					// if the task is blocking, it will be handled by a worker from the pool
					workerPool <- struct{}{} // aquire a worker from the pool (wait if none available)

					// Execute the blocking task in a separate goroutine
					go func() {
						defer func() {
							<-workerPool // release the worker back to the pool once the task is completed
						}()

						// run the main task
						task.MainTask()

						// if a callback is provided, add it to the taskQueue
						if task.Callback != nil {
							eventLoop.AddTaskToQueue(&Task{
								MainTask: task.Callback, // the callback will be executed after the main task completes
							})
						}
					}()
				} else {
					// for non blocking tasks, execute immediately in the event loop
					task.MainTask()
				}
			case task := <-eventLoop.taskQueue:
				// Execute the callback task
				task.MainTask()

				// Handle the stop signal to terminate the event loop
			case stop := <-eventLoop.stop:
				if stop {
					return // Exit the loop and terminate the event loop
				}
			}
		}
	}()
	return &wg
}
