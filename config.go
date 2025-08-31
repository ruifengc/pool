package pool

import "log"

type PoolConfig struct {
	WorkerCount  int
	TaskChanSize int
	Logger       Logger
}

var DefaultPoolConfig = &PoolConfig{
	WorkerCount:  8,
	TaskChanSize: 4,
	Logger:       log.Default(),
}

type Option func(*PoolConfig)

// WithWorkerCount sets the number of workers.
func WithWorkerCount(count int) Option {
	return func(p *PoolConfig) {
		p.WorkerCount = count
	}
}

// WithTaskChanSize sets the size of the task channel.
func WithTaskChanSize(size int) Option {
	return func(p *PoolConfig) {
		p.TaskChanSize = size
	}
}

// WithLogger sets the logger.
func WithLogger(logger Logger) Option {
	return func(p *PoolConfig) {
		p.Logger = logger
	}
}
