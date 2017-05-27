package main

import (
	"fmt"
	"task"
	"time"
)

func main() {
	apps := []*task.App{}

	// Task specifications
	taskSpecs := []task.TaskSpec{
		task.TaskSpec{
			Period:           4 * time.Second,
			TotalRunTimeMean: 1 * time.Second,
			TotalRunTimeStd:  300 * time.Millisecond,
			RelativeDeadline: 3 * time.Second,
		},
		task.TaskSpec{
			Period:           2 * time.Second,
			TotalRunTimeMean: 1 * time.Second,
			TotalRunTimeStd:  300 * time.Millisecond,
			RelativeDeadline: 2 * time.Second,
		},
	}

	// Create all applications
	for i, taskSpec := range taskSpecs {
		apps = append(apps, task.NewApp(fmt.Sprintf("app%d", i), taskSpec))
	}

	/*
		// Create and initialize the scheduler
		sched := scheduler.NewScheduler()
		// To be implemented, initialization process

		// Start the scheduler
		sched.Start()

		// Start all applications
		for _, app := range apps {
			app.Start()
		}

		time.Sleep(c.TEST_TIME)

		// Stop all applications
		for _, app := range apps {
			app.Stop()
		}

		// Stop the scheduler
		sched.Stop()
	*/
}
