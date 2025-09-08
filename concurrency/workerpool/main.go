package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Job represents work to be done
type Job struct {
	ID       int
	Duration time.Duration
	Result   chan string
}

// Worker represents a worker that processes jobs
type Worker struct {
	ID   int
	Jobs chan Job
	Quit chan bool
}

// WorkerPool manages a pool of workers
type WorkerPool struct {
	Workers    []*Worker
	JobQueue   chan Job
	ResultChan chan string
	wg         sync.WaitGroup
}

// NewWorker creates a new worker
func NewWorker(id int, jobQueue chan Job) *Worker {
	return &Worker{
		ID:   id,
		Jobs: jobQueue,
		Quit: make(chan bool),
	}
}

// Start starts the worker
func (w *Worker) Start(wg *sync.WaitGroup, resultChan chan string) {
	defer wg.Done()
	go func() {
		for {
			select {
			case job := <-w.Jobs:
				fmt.Printf("Worker %d started job %d\n", w.ID, job.ID)
				
				// Simulate work
				time.Sleep(job.Duration)
				
				result := fmt.Sprintf("Job %d completed by worker %d", job.ID, w.ID)
				resultChan <- result
				job.Result <- result
				close(job.Result)
				
				fmt.Printf("Worker %d finished job %d\n", w.ID, job.ID)
				
			case <-w.Quit:
				fmt.Printf("Worker %d stopping\n", w.ID)
				return
			}
		}
	}()
}

// Stop stops the worker
func (w *Worker) Stop() {
	go func() {
		w.Quit <- true
	}()
}

// NewWorkerPool creates a new worker pool
func NewWorkerPool(numWorkers int) *WorkerPool {
	jobQueue := make(chan Job, 100)
	resultChan := make(chan string, 100)
	
	pool := &WorkerPool{
		Workers:    make([]*Worker, numWorkers),
		JobQueue:   jobQueue,
		ResultChan: resultChan,
	}
	
	// Create workers
	for i := 0; i < numWorkers; i++ {
		worker := NewWorker(i+1, jobQueue)
		pool.Workers[i] = worker
	}
	
	return pool
}

// Start starts all workers in the pool
func (p *WorkerPool) Start() {
	for _, worker := range p.Workers {
		p.wg.Add(1)
		worker.Start(&p.wg, p.ResultChan)
	}
}

// AddJob adds a job to the pool
func (p *WorkerPool) AddJob(job Job) {
	p.JobQueue <- job
}

// Stop stops all workers
func (p *WorkerPool) Stop() {
	for _, worker := range p.Workers {
		worker.Stop()
	}
	p.wg.Wait()
	close(p.JobQueue)
	close(p.ResultChan)
}

func main() {
	// Create worker pool with 3 workers
	pool := NewWorkerPool(3)
	pool.Start()
	
	// Start result collector
	go func() {
		for result := range pool.ResultChan {
			fmt.Printf("✓ %s\n", result)
		}
	}()
	
	// Generate random jobs
	rand.Seed(time.Now().UnixNano())
	numJobs := 10
	
	var jobResults []chan string
	
	for i := 1; i <= numJobs; i++ {
		resultChan := make(chan string, 1)
		jobResults = append(jobResults, resultChan)
		
		job := Job{
			ID:       i,
			Duration: time.Duration(rand.Intn(3)+1) * time.Second,
			Result:   resultChan,
		}
		
		pool.AddJob(job)
		fmt.Printf("Added job %d\n", i)
	}
	
	// Wait for all jobs to complete
	fmt.Println("Waiting for all jobs to complete...")
	for i, resultChan := range jobResults {
		result := <-resultChan
		fmt.Printf("Job %d result: %s\n", i+1, result)
	}
	
	// Stop the pool
	fmt.Println("Stopping worker pool...")
	pool.Stop()
	
	fmt.Println("All jobs completed!")
}
