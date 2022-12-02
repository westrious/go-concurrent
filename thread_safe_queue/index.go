package thread_safe_queue

import "sync"

//ThreadSafeQueueInt64 线程安全队列
type ThreadSafeQueueInt64 struct {
	queue []int64
	mu    sync.Mutex
}

func NewThreadSafeQueueInt64(n int) *ThreadSafeQueueInt64 {
	return &ThreadSafeQueueInt64{queue: make([]int64, 0, n)}
}

func (t *ThreadSafeQueueInt64) Enqueue(v int64) {
	t.mu.Lock()
	t.queue = append(t.queue, v)
	t.mu.Unlock()
}

func (t *ThreadSafeQueueInt64) Dequeue() int64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	if len(t.queue) <= 0 {
		return 0
	}
	v := t.queue[0]
	t.queue = t.queue[1:]
	return v
}
