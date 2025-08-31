package pool

import (
	"log"
	"sync"
	"time"
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
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Worker panic recovered: %v", r)
					}
					wg.Done()
				}()

				task()
			}()
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

// AddTaskWithTimeout 添加带超时控制的任务
func (p *Pool) AddTaskWithTimeout(task Task, timeout time.Duration) {
	p.wg.Add(1)
	p.taskChan <- func() {
		done := make(chan struct{})

		go func() {
			defer close(done)
			task()
		}()
		timer := time.NewTimer(timeout)
		defer timer.Stop()

		select {
		case <-done:
			// 任务正常完成
		case <-timer.C:
			// 超时处理，记录日志而不是 panic
			log.Printf("Task timeout after %v", timeout)
		}
	}
}

func (p *Pool) Wait() {
	p.wg.Wait()
}

func (p *Pool) Release() {
	close(p.taskChan)
	p.Wait()
}
