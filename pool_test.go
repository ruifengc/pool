package pool

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	// 创建一个有 3 个工作协程的协程池
	pool := NewPool(WithWorkerCount(3))

	// 关闭协程池
	defer pool.Release()

	var counter int32

	// 添加任务
	for i := 0; i < 100; i++ {
		pool.AddTask(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	// 等待所有任务完成
	pool.Wait()

	if counter != 100 {
		t.Errorf("Expected counter to be 100, but got %d", counter)
	}
}

func TestPoolWithTimeout_NormalTask(t *testing.T) {
	pool := NewPool(WithWorkerCount(3), WithTaskChanSize(100))
	defer pool.Release()

	var counter int32

	// 添加正常任务（不会超时）
	pool.AddTaskWithTimeout(func() {
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	}, 200*time.Millisecond)

	pool.Wait()

	if counter != 1 {
		t.Errorf("Expected counter to be 1, but got %d", counter)
	}
}

func TestPoolWithTimeout_TimeoutTask(t *testing.T) {
	pool := NewPool(WithWorkerCount(3))
	defer pool.Release()

	var counter int32

	// 添加会超时的任务
	pool.AddTaskWithTimeout(func() {
		time.Sleep(500 * time.Millisecond) // 超过超时时间
		atomic.AddInt32(&counter, 1)
	}, 100*time.Millisecond)

	pool.Wait()

	// 超时任务应该不会增加计数器
	if counter != 0 {
		t.Errorf("Expected counter to be 0 due to timeout, but got %d", counter)
	}
}

func TestPoolWithTimeout_MixedTasks(t *testing.T) {
	pool := NewPool(WithWorkerCount(3))
	defer pool.Release()

	var normalCounter int32
	var timeoutCounter int32

	// 添加正常任务
	for i := 0; i < 5; i++ {
		pool.AddTaskWithTimeout(func() {
			time.Sleep(20 * time.Millisecond)
			atomic.AddInt32(&normalCounter, 1)
		}, 100*time.Millisecond)
	}

	// 添加超时任务
	for i := 0; i < 3; i++ {
		pool.AddTaskWithTimeout(func() {
			time.Sleep(200 * time.Millisecond) // 超过超时时间
			atomic.AddInt32(&timeoutCounter, 1)
		}, 50*time.Millisecond)
	}

	pool.Wait()

	if normalCounter != 5 {
		t.Errorf("Expected normalCounter to be 5, but got %d", normalCounter)
	}

	if timeoutCounter != 0 {
		t.Errorf("Expected timeoutCounter to be 0, but got %d", timeoutCounter)
	}
}

func TestPoolWithTimeout_ImmediateTask(t *testing.T) {
	pool := NewPool(WithTaskChanSize(3))
	defer pool.Release()

	var counter int32

	// 添加立即完成的任务
	pool.AddTaskWithTimeout(func() {
		atomic.AddInt32(&counter, 1)
	}, time.Second)

	pool.Wait()

	if counter != 1 {
		t.Errorf("Expected counter to be 1, but got %d", counter)
	}
}

func TestPoolWithTimeout_ZeroTimeout(t *testing.T) {
	pool := NewPool(WithTaskChanSize(3))
	defer pool.Release()

	var counter int32

	// 添加零超时的任务
	pool.AddTaskWithTimeout(func() {
		time.Sleep(10 * time.Millisecond)
		atomic.AddInt32(&counter, 1)
	}, 0) // 零超时应该立即超时

	pool.Wait()

	if counter != 0 {
		t.Errorf("Expected counter to be 0 with zero timeout, but got %d", counter)
	}
}

func TestPoolConcurrency(t *testing.T) {
	pool := NewPool(WithTaskChanSize(100), WithWorkerCount(30))
	defer pool.Release()

	var counter int32

	// 并发添加大量任务
	for i := 0; i < 1000; i++ {
		pool.AddTask(func() {
			atomic.AddInt32(&counter, 1)
		})
	}

	pool.Wait()

	if counter != 1000 {
		t.Errorf("Expected counter to be 1000, but got %d", counter)
	}
}
