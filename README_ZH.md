# Pool

[English](README.md) | [中文](README_ZH.md)

这是 Go 中 coroutine（goroutine）池的一个简单实现。该池可控制并发 goroutines 的数量，从而让您在有限并发的情况下高效地管理和执行大量任务。
## 特性

- 管理 goroutines 池以控制并发性。
- 向池中添加任务，以便并发执行。
- 等待所有任务完成后再关闭池。
- 所有任务完成后，自动停止工作进程。

## 用例
下面是一个如何使用 coroutine 池的基本示例：

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