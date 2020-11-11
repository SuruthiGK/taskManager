package taskManager

import (
	"fmt"
)

type Job struct {
  ID int
  IsCompleted bool	// have a random function to mark the IsCompleted after a random period  
  Status string			// completed, failed, timeout
  //Time string			// when was the task created
  TaskData string		// random string indicating task
}


type Worker struct {
	ID int
	jobs chan *Job
	dispatchStatus chan *DispatchStatus
	Quit chan bool
}
 
func CreateNewWorker(id int, workerQueue chan *Worker, jobQueue chan *Job, dStatus chan *DispatchStatus) *Worker {
	w := &Worker{
		ID: id, 
		jobs: jobQueue,
		dispatchStatus: dStatus,
	}

	go func() { workerQueue <- w }()
	return w
}
 
func (w *Worker) Start() {
	go func() {
		for {
			select {
			case job := <- w.jobs:
				fmt.Println("Worker[%d] executing job[%d].\n", w.ID, job.ID)
				//w.dispatchStatus <- &DispatchStatus{Type: "worker", ID: w.ID, Status: "quit"}
				w.dispatchStatus <- &DispatchStatus{Type: "worker", ID: w.ID, Status: job.Status}
				w.Quit <- true
			case <- w.Quit:
				return
			}
		}
	}()
}

type DispatchStatus struct {
	Type string
	ID int
	Status string
}
 
type Dispatcher struct {
	jobCounter int                      // internal counter for number of jobs
	jobQueue chan *Job                  // channel of jobs submitted by main()
	dispatchStatus chan *DispatchStatus // channel for job/worker status reports
	workQueue chan *Job                 // channel of work dispatched
	workerQueue chan *Worker            // channel of workers
}
func CreateNewDispatcher() *Dispatcher {
	d := &Dispatcher{
		jobCounter: 0,
		jobQueue: make(chan *Job),
		dispatchStatus: make(chan *DispatchStatus),
		workQueue: make(chan *Job),
		workerQueue: make(chan *Worker),
	}
return d
}

//type JobExecutable func() error
 
func (d *Dispatcher) Start(numWorkers int) {
	// Create numWorkers:
	for i := 0; i<numWorkers; i++ {
		worker := CreateNewWorker(i, d.workerQueue, d.workQueue, d.dispatchStatus)
		worker.Start()
	}

	// wait for work to be added then pass it off.
	go func() { 
		for {
			select {
			case job := <- d.jobQueue:
				fmt.Println("Got a job in the queue to dispatch: %d\n", job.ID)
				// Sending it off;
				//fmt.Println(job)
				
				if job.Status != "Completed" {
					fmt.Println(job)
					d.workQueue <- job
				}else {
					d.jobCounter--
				}
								
			case ds := <- d.dispatchStatus:
				fmt.Println("Got a dispatch status:\n\tType[%s] - ID[%d] - Status[%s]\n", ds.Type, ds.ID, ds.Status)
				
				if ds.Type == "worker" {
					if ds.Status == "timeout" {
						d.jobCounter--
					}
				}
			}
		}
}()
}

func (d *Dispatcher) AddJob(je Job) {
	j := &Job{ID: d.jobCounter, IsCompleted: je.IsCompleted, Status: je.Status,	TaskData: je.TaskData}
	go func() { d.jobQueue <- j }()
	d.jobCounter++
	fmt.Println("jobCounter is now: %d\n", d.jobCounter)
}

func (d *Dispatcher) Finished() bool {
	if d.jobCounter < 1 {
		return true
	} else {
		return false
	}
}

func Work() {

	fmt.Println("job1: performing work.\n")
	job1 := Job{IsCompleted: false,Status: "In-Progress",TaskData: "email jack"}

	job2 := Job{IsCompleted: true,Status: "Completed",TaskData: "email Rose"}

	d := CreateNewDispatcher()
	d.AddJob(job1)
	d.AddJob(job2)
	d.Start(2)

	for {
		if d.Finished() {
			fmt.Println("All jobs finished.\n")
			break
		}
	}
}