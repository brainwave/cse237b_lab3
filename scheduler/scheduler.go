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

//NewScheduler assigns global channels created in main.go
func NewScheduler(t chan *task.Task, w chan *worker.Worker) *Scheduler {
	log.Printf("Creating new scheduler\n")

	//Create a dummy scheduler object
	sched := Scheduler{
		TaskChan:      t,
		WorkerChan:    w,
		StopChan:      make(chan interface{}),
		FreeWorkerBuf: &worker.WorkerPool{},
		AllWorkerBuf:  &worker.WorkerPool{},
		TaskBuf:       &worker.TaskQueue{},
	}

	log.Println("Done starting scheduler")
	return &sched
}

// ScheduleLoop runs the scheduling algorithm inside a goroutine
func (s *Scheduler) ScheduleLoop() {
	log.Printf("Scheduler: Scheduling loop starts\n")
loop:
	for {
		select {
		case newTask := <-s.TaskChan:
			// Receive a new task and do scheduling
			s.TaskBuf.Queue = append(s.TaskBuf.Queue, newTask)

			TaskLen := len(s.TaskBuf.Queue) - 1
			WrkrLen := len(s.FreeWorkerBuf.Pool) - 1

			log.Printf("Scheduler: Task Recieved, task%d %s, inserted. Queue Len: %d\n", newTask.TaskID, newTask.AppID, TaskLen+1)

			if TaskLen >= 0 && WrkrLen >= 0 {
				log.Printf("Scheduler: Task Handover, task%d %s -> %d\n", s.TaskBuf.Queue[TaskLen].TaskID, s.TaskBuf.Queue[TaskLen].AppID, s.FreeWorkerBuf.Pool[WrkrLen].WorkerID)

				s.FreeWorkerBuf.Pool[WrkrLen].TaskChan <- s.TaskBuf.Queue[TaskLen]
				s.FreeWorkerBuf.Pool = s.FreeWorkerBuf.Pool[:WrkrLen]
				s.TaskBuf.Queue = s.TaskBuf.Queue[:TaskLen]

				//(1) indicates that the handover was caused by new task
				log.Printf("Scheduler: Task Handover Complete(1)\n")

			}

		case w := <-s.WorkerChan:
			// Assign the worker a new task, and remove the task from queue
			TaskLen := len(s.TaskBuf.Queue) - 1
			WrkrLen := len(s.FreeWorkerBuf.Pool) - 1

			s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, w)
			log.Printf("Scheduler: Worker %d free, inserted. num free workers: %d", w.WorkerID, WrkrLen+1)

			if TaskLen >= 0 && WrkrLen >= 0 {

				log.Printf("Scheduler: Task Handover: task%d, %s -> Worker %d\n", s.TaskBuf.Queue[TaskLen].TaskID,
					s.TaskBuf.Queue[TaskLen].AppID, w.WorkerID)

				s.FreeWorkerBuf.Pool[WrkrLen-1].TaskChan <- s.TaskBuf.Queue[TaskLen]
				s.FreeWorkerBuf.Pool = s.FreeWorkerBuf.Pool[:WrkrLen]
				s.TaskBuf.Queue = s.TaskBuf.Queue[:TaskLen]

				//(2) indicates that the handover was caused by worker becoming free
				log.Printf("Scheduler: Task Handover Complete(2)\n")
			}

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
	i := 0

	go s.ScheduleLoop()

	s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, new(worker.Worker))
	s.FreeWorkerBuf.Pool = append(s.FreeWorkerBuf.Pool, new(worker.Worker))

	for _, wrkr := range s.FreeWorkerBuf.Pool {
		wrkr.WorkerID = i
		wrkr.TaskChan = make(chan *task.Task)
		wrkr.WorkerChan = s.WorkerChan

		go wrkr.TaskProcessLoop()
		i = i + 1
	}

	s.AllWorkerBuf = s.FreeWorkerBuf

}

// Stop stops the scheduler
func (s *Scheduler) Stop() {
	s.StopChan <- 0
	for _, w := range s.AllWorkerBuf.Pool {
		w.StopChan <- 0
	}
}
