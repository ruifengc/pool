# Pool

[English](README.md) | [中文](README_ZH.md)

这是一个高性能的Go语言goroutine池实现，提供了并发控制、任务管理和超时处理等功能。

## 特性

- 🚀 **高性能并发控制** - 精确控制goroutine并发数量
- ⏰ **超时任务支持** - 内置任务超时处理机制
- 🔧 **灵活配置** - 支持工作线程数、任务通道大小、日志器等配置
- 🛡️ **异常恢复** - 自动捕获并处理worker panic
- 📊 **任务统计** - 支持等待所有任务完成

## 安装

```bash
go get github.com/ruifengc/pool
```

## 快速开始

### 基础用法

```go
package main

import (
    "fmt"
    "github.com/ruifengc/pool"
)

func main() {
    // 创建一个包含3个worker的池
    pool := pool.NewPool(pool.WithWorkerCount(3))
	var counter int32

    defer func(){
		pool.Release()
		// Release()自动调用Wait()
		fmt.Printf("All tasks completed, counter: %d\n", counter)
	}()

    // 添加任务到池中
    for i := 0; i < 10; i++ {
        taskNum := i
        pool.AddTask(func() {
            fmt.Printf("Executing task %d\n", taskNum)
            atomic.AddInt32(&counter, 1)
        })
    }
}
```

### 带超时的任务

```go
func main() {
    pool := pool.NewPool(pool.WithWorkerCount(2))
    defer pool.Release()

    // 添加带超时的任务
    pool.AddTaskWithTimeout(func() {
        time.Sleep(2 * time.Second)
        fmt.Println("这个任务会正常完成")
    }, 3*time.Second)

    // 添加会超时的任务
    pool.AddTaskWithTimeout(func() {
        time.Sleep(5 * time.Second) // 超过超时时间
        fmt.Println("这个任务不会完成")
    }, 1*time.Second)

    pool.Wait()
}
```

## 配置选项

### 工作线程数

```go
// 设置10个工作线程
pool := pool.NewPool(pool.WithWorkerCount(10))
```

### 任务通道大小

```go
// 设置任务通道大小为100
pool := pool.NewPool(pool.WithTaskChanSize(100))
```

### 自定义日志器

```go
type MyLogger struct{}

func (l *MyLogger) Printf(format string, v ...any) {
    log.Printf("[POOL] "+format, v...)
}

pool := pool.NewPool(pool.WithLogger(&MyLogger{}))
```

## API 参考

### NewPool(options ...Option) *Pool

创建新的goroutine池实例。

### AddTask(task Task)

向池中添加任务。

### AddTaskWithTimeout(task Task, timeout time.Duration)

添加带超时控制的任务。

### Wait()

等待所有任务完成。

### Release()

释放池资源，关闭任务通道。

## 性能建议

1. **工作线程数**：通常设置为CPU核心数的2-4倍
2. **任务通道大小**：根据任务产生速度调整，避免阻塞
3. **超时设置**：合理设置超时时间，避免资源浪费

## 测试

运行测试套件：

```bash
go test -v ./...
```

## 贡献

欢迎提交Issue和Pull Request来改进这个项目。