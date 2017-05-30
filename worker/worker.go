package worker

import (
	"log"
	"sync"
	"task"
	"time"
)

type TaskQueue struct {
	Queue []*task.Task
	Lock  *sync.Mutex
}

type WorkerPool struct {
	Pool []*Worker
	Lock *sync.Mutex
}

// Worker is the agent to process tasks
type Worker struct {
	WorkerID   int
	TaskChan   chan *task.Task
	WorkerChan chan *Worker
	StopChan   chan interface{}
}

// TaskProcessLoop processes tasks without preemption
func (w *Worker) TaskProcessLoop() {
	log.Printf("Worker<%d>: Task processor starts\n", w.WorkerID)
loop:
	for {
		select {
		case t := <-w.TaskChan:
			// This worker receives a new task to run
			w.Process(t)
			log.Printf("Worker <%d>: App<%s>/Task<%d> ends\n", w.WorkerID, t.AppID, t.TaskID)
			log.Printf("Conveying worker %d state to channel", w.WorkerID)
			w.WorkerChan <- w

		case <-w.StopChan:
			// Receive signal to stop
			// To be implemented
			break loop
		}
	}
	log.Printf("Worker<%d>: Task processor ends\n", w.WorkerID)
}

// Process runs a task on a worker without preemption
func (w *Worker) Process(t *task.Task) {
	// Process the task
	time.Sleep(t.TotalRunTime)
	log.Printf("Worker <%d>: App<%s>/Task<%d> ends\n", w.WorkerID, t.AppID, t.TaskID)
}

func (tq TaskQueue) Less(i, j int) bool {
	return (tq.Queue[i].Deadline.Before(tq.Queue[j].Deadline))
}

func (tq TaskQueue) Swap(i, j int) {
	tq.Queue[i], tq.Queue[j] = tq.Queue[j], tq.Queue[i]
}

func (tq TaskQueue) Len() int {
	return (len(tq.Queue))
}
