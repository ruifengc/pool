package pool

import (
	"sync"
)

type Task func()

type Worker struct {
	taskChan chan Task
}

func NewWorker(taskChan chan Task) Worker {
	return Worker{
		taskChan: taskChan,
	}
}

func (w Worker) Start(wg *sync.WaitGroup) {
	go func() {
		for task := range w.taskChan {
			task()
			wg.Done()
		}
	}()
}

type Pool struct {
	taskChan chan Task
	workers  []Worker
	wg       *sync.WaitGroup
}

func NewPool(workerCount int) *Pool {
	taskChan := make(chan Task)
	workers := make([]Worker, workerCount)
	wg := &sync.WaitGroup{}

	pool := &Pool{
		taskChan: taskChan,
		workers:  workers,
		wg:       wg,
	}

	for i := 0; i < workerCount; i++ {
		worker := NewWorker(taskChan)
		worker.Start(wg)
		pool.workers[i] = worker
	}

	return pool
}

func (p *Pool) AddTask(task Task) {
	p.wg.Add(1)
	p.taskChan <- task
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	close(p.taskChan)
}
