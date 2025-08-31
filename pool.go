package pool

import (
	"sync"
	"time"
)

type Task func()

type Worker struct {
	pool *Pool
}

func NewWorker(pool *Pool) Worker {
	return Worker{
		pool: pool,
	}
}

func (w Worker) Start(wg *sync.WaitGroup) {
	go func() {
		for task := range w.pool.taskChan {
			func() {
				defer func() {
					if r := recover(); r != nil {
						w.pool.log.Printf("Worker panic recovered: %v", r)
					}
					wg.Done()
				}()

				task()
			}()
		}
	}()
}

type Logger interface {
	Printf(format string, v ...any)
}

// pool
type Pool struct {
	taskChan chan Task
	workers  []Worker
	// logger
	log Logger
	// wait group
	wg *sync.WaitGroup
}

func NewPool(Options ...Option) *Pool {
	config := DefaultPoolConfig
	for _, opt := range Options {
		opt(config)
	}

	taskChan := make(chan Task, config.TaskChanSize)
	workers := make([]Worker, config.WorkerCount)

	wg := &sync.WaitGroup{}

	pool := &Pool{
		taskChan: taskChan,
		workers:  workers,
		wg:       wg,
		log:      config.Logger,
	}

	for i := 0; i < config.WorkerCount; i++ {
		worker := NewWorker(pool)
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
			p.log.Printf("Task timeout after %v", timeout)
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
