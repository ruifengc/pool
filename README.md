# Pool

English | [中文](README_ZH.md)

This is a simple implementation of a coroutine (goroutine) pool in Go. The pool controls the number of concurrent goroutines, allowing you to efficiently manage and execute a large number of tasks with limited concurrency.

## Features

- Manage a pool of goroutines to control concurrency.
- Add tasks to the pool for concurrent execution.
- Wait for all tasks to complete before shutting down the pool.
- Automatically stop workers once all tasks are done.

## Usage
  Here is a basic example of how to use the coroutine pool:

```go
package main

import (
	"fmt"
	"github.com/ruifengc/pool"
)

func main() {
	// Create a pool with 3 workers
	pool := NewPool(3)
	// Release the pool resources
    defer pool.Release()
	// Add tasks to the pool
	for i := 0; i < 10; i++ {
		taskNum := i
		pool.AddTask(func() {
			fmt.Printf("Executing task %d\n", taskNum)
		})
	}

	// Wait for all tasks to complete
	pool.Wait()

	fmt.Println("All tasks completed")
}

```