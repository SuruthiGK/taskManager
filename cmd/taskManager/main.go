package main

import (
	"fmt"
	"taskManager/internal/taskManager"
)

// Start all the listeners, and reconnect and wait for them to be done.
// Route the data from listeners to the right channels.
// Configure the worker count for each of these channels.
// Worker by definition will have to call a function based on the message.

func main() {
	fmt.Println("started...\n")
	taskManager.Work()
}
