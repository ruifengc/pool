# Pool

[English](README.md) | [中文](README_ZH.md)

A high-performance goroutine pool implementation in Go with concurrency control, task management, and timeout handling.

## Features

- 🚀 **High-performance concurrency control** - Precise control over goroutine concurrency
- ⏰ **Timeout task support** - Built-in task timeout handling
- 🔧 **Flexible configuration** - Configurable worker count, task channel size, and logger
- 🛡️ **Exception recovery** - Automatic panic recovery for worker goroutines
- 📊 **Task statistics** - Wait for all tasks to complete

## Installation

```bash
go get github.com/ruifengc/pool
```

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "sync/atomic"
    "github.com/ruifengc/pool"
)

func main() {
    // Create a pool with 3 workers
    pool := pool.NewPool(pool.WithWorkerCount(3))
	var counter int32

    defer func(){
		pool.Release()
		// Release() automatically calls Wait()
		fmt.Printf("All tasks completed, counter: %d\n", counter)
	}()

    // Add tasks to the pool
    for i := 0; i < 10; i++ {
        taskNum := i
        pool.AddTask(func() {
            fmt.Printf("Executing task %d\n", taskNum)
            atomic.AddInt32(&counter, 1)
        })
    }
}
```

### Timeout Tasks

```go
func main() {
    pool := pool.NewPool(pool.WithWorkerCount(2))
    defer pool.Release()

    // Add task with timeout (will complete normally)
    pool.AddTaskWithTimeout(func() {
        time.Sleep(2 * time.Second)
        fmt.Println("This task will complete normally")
    }, 3*time.Second)

    // Add task that will timeout
    pool.AddTaskWithTimeout(func() {
        time.Sleep(5 * time.Second) // Exceeds timeout
        fmt.Println("This task won't complete")
    }, 1*time.Second)

    pool.Wait()
}
```

## Configuration Options

### Worker Count

```go
// Set 10 workers
pool := pool.NewPool(pool.WithWorkerCount(10))
```

### Task Channel Size

```go
// Set task channel size to 100
pool := pool.NewPool(pool.WithTaskChanSize(100))
```

### Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Printf(format string, v ...any) {
    log.Printf("[POOL] "+format, v...)
}

pool := pool.NewPool(pool.WithLogger(&MyLogger{}))
```

## API Reference

### NewPool(options ...Option) *Pool

Creates a new goroutine pool instance.

### AddTask(task Task)

Adds a task to the pool.

### AddTaskWithTimeout(task Task, timeout time.Duration)

Adds a task with timeout control.

### Wait()

Waits for all tasks to complete.

### Release()

Releases pool resources and closes the task channel.

## Performance Tips

1. **Worker count**: Typically 2-4 times the number of CPU cores
2. **Task channel size**: Adjust based on task generation rate to avoid blocking
3. **Timeout settings**: Set reasonable timeout durations to avoid resource waste

## Testing

Run the test suite:

```bash
go test -v ./...
```

## Contributing

Issues and Pull Requests are welcome to improve this project.