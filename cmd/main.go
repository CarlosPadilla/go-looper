package main

import (
	"fmt"
	"looper/pkg/eventloop"
	"time"
)

func main() {
	// Initialize the event loop
	eventL := eventloop.NewEventLoop(10)
	// Initialize the event loop and get the wait group to track the completion of tasks
	wg := eventL.InitEventLoop(10)
	// Push tasks to the event loop
	eventL.Add(&eventloop.Task{
		MainTask: func() {
			fmt.Println("Non-blocking task 1 executed.")
		},
		IsBlocking: false,
	})
	eventL.Add(&eventloop.Task{
		MainTask: func() {
			fmt.Println("Non-blocking task 2 executed.")
		},
		IsBlocking: false,
	})
	// Add blocking tasks with callbacks
	eventL.Add(&eventloop.Task{
		MainTask: func() {
			time.Sleep(0 * time.Second) // Simulate a blocking task
			fmt.Println("Blocking task 3 executed.")
		},
		Callback: func() {
			fmt.Println("Callback to blocking task 3 executed.")
		},
		IsBlocking: true,
	})
	eventL.Add(&eventloop.Task{
		MainTask: func() {
			time.Sleep(4 * time.Second) // Simulate a blocking task
			fmt.Println("Blocking task 4 executed.")
		},
		Callback: func() {
			fmt.Println("Callback to blocking task 4 executed.")
		},
		IsBlocking: true,
	})
	// Push another non-blocking task
	eventL.Add(&eventloop.Task{
		MainTask: func() {
			fmt.Println("Non-blocking task 5 executed.")
		},
		IsBlocking: false,
	})
	// Simulate a wait to let all the operations to finish
	time.Sleep(10 * time.Second)

	// Indicate the event loop to stop running
	eventL.StopEventLoop()
	// Wait for all tasks to finish before exiting the program
	wg.Wait()
}
