package my_mutex

import (
	"sync"
	"sync/atomic"
	"time"
	"unsafe"
)

// 复制Mutex定义的常量
const (
	mutexLocked      = 1 << iota // 加锁标识位置
	mutexWoken                   // 唤醒标识位置（标记是否有“通过unlock唤醒”的waiter在竞争锁）
	mutexStarving                // 锁饥饿标识位置
	mutexWaiterShift = iota      // 标识waiter的起始bit位置
)

type Mutex struct {
	sync.Mutex
}

func (m *Mutex) TryLock() bool {
	// 如果当前锁处于空闲
	if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), 0, mutexLocked) {
		return true
	}

	// 如果处于唤醒、加锁或者饥饿状态，这次请求就不参与竞争了，返回 false
	old := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	if old&(mutexLocked|mutexStarving|mutexWoken) != 0 {
		return false
	}

	// 尝试在竞争的情况下请求锁
	return atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(&m.Mutex)), old, old|mutexLocked)
}

//GetWaiterNum 获取锁的等待队列中的 goroutine 数量
func (m *Mutex) GetWaiterNum() int32 {
	// 获取 state 字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	v = v >> mutexWaiterShift
	return v
}

//Count 获取抢占锁的 goroutine 数量
func (m *Mutex) Count() int32 {
	// 获取 state 字段的值
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	// 等待者 + 拥有者
	v = v>>mutexWaiterShift + (v & mutexLocked)
	return v
}

//IsLocked 锁是否被持有
func (m *Mutex) IsLocked() bool {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return v&mutexLocked == mutexLocked
}

//IsWoken 是否有等待者被唤醒
func (m *Mutex) IsWoken() bool {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return v&mutexWoken == mutexWoken
}

//IsStarving 锁是否处于饥饿状态
func (m *Mutex) IsStarving() bool {
	v := atomic.LoadInt32((*int32)(unsafe.Pointer(&m.Mutex)))
	return v&mutexStarving == mutexStarving
}

//LockWithTimeout 有超时时间的 lock
func (m *Mutex) LockWithTimeout(timeout time.Duration) bool {
	ch := make(chan bool)
	go func() {
		for {
			select {
			case <-time.After(timeout):
				ch <- false
				return
			default:
				if m.TryLock() {
					ch <- true
					return
				}
			}
		}
	}()
	return <-ch
}
