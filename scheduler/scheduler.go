package scheduler

import (
	"log"
	"task"
	"worker"
)

// Scheduler dispatches tasks to workers
type Scheduler struct {
	TaskChan      chan *task.Task
	WorkerChan    chan *worker.Worker
	StopChan      chan interface{}
	FreeWorkerBuf *worker.WorkerPool
	AllWorkerBuf  *worker.WorkerPool
	TaskBuf       *worker.TaskQueue
}

func NewScheduler(t chan *task.Task, w chan *worker.Worker) *Scheduler {
	log.Printf("Creating new scheduler\n")

	return &Scheduler{
		TaskChan:   t,
		WorkerChan: w,
		StopChan:   make(chan interface{}),
	}
}

// ScheduleLoop runs the scheduling algorithm inside a goroutine
func (s *Scheduler) ScheduleLoop() {
	log.Printf("Scheduler: Scheduling loop starts\n")
loop:
	for {
		select {
		case newTask := <-s.TaskChan:
			// Receive a new task and do scheduling
			log.Printf("Scheduler: New Task, %d\n", newTask.TaskID)
		case w := <-s.WorkerChan:
			// A worker becomes free
			log.Printf("Scheduler: Worker Free, %d\n", w.WorkerID)
		case <-s.StopChan:
			// Receive signal to stop scheduling
			log.Printf("Scheduler: Stop Signal\n")
			break loop
		}
	}
	log.Printf("Scheduler: Task processor ends\n")
}

// Start starts the scheduler
func (s *Scheduler) Start() {
	go s.ScheduleLoop()
}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.StopChan <- 0
	for _, w := range s.AllWorkerBuf.Pool {
		w.StopChan <- 0
	}
}
