package my_rw_mutex

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const rwmutexMaxReaders = 1 << 30

type RWMutex struct {
	sync.RWMutex
}

//TryRLock 尝试获取读锁
func (m *RWMutex) TryRLock() bool {
	for {
		readerCount := atomic.LoadInt32((*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&m.RWMutex)) + unsafe.Sizeof(sync.Mutex{}) + 2*unsafe.Sizeof(uint32(0)))))
		// 有 writer 获取锁，直接返回返回 false
		if readerCount < 0 {
			return false
		}
		// 尝试获取读锁
		if atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&m.RWMutex))+unsafe.Sizeof(sync.Mutex{})+2*unsafe.Sizeof(uint32(0)))), readerCount, readerCount+1) {
			return true
		}
	}
}

//TryLock 尝试获取写锁
func (m *RWMutex) TryLock() bool {
	// 这里不想再写一遍 mutex 的 tryLock
	if !m.RWMutex.w.TryLock() {
		return false
	}
	if !atomic.CompareAndSwapInt32((*int32)(unsafe.Pointer(uintptr(unsafe.Pointer(&m.RWMutex))+unsafe.Sizeof(sync.Mutex{})+2*unsafe.Sizeof(uint32(0)))), 0, -rwmutexMaxReaders) {
		// 这里不想再写一遍 mutex 的 UnLock
		m.RWMutex.w.Unlock()
		return false
	}
	return true
}
