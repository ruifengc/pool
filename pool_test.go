package pool

import (
	"sync/atomic"
	"testing"
)

func TestPool(t *testing.T) {
	// 创建一个有 3 个工作协程的协程池
	pool := NewPool(3)

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
