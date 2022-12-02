package recursive_mutex

import (
	"fmt"
	"go-concurrent/utils"
	"sync"
	"sync/atomic"
)

// RecursiveMutex 可重入锁
type RecursiveMutex struct {
	sync.Mutex
	owner     int64 // 当前持有锁的 goroutine id
	recursion int32 // 这个 goroutine 重入的次数
}

func (m *RecursiveMutex) Lock() {
	gid := utils.GoID()
	if atomic.LoadInt64(&m.owner) == gid {
		m.recursion++
		return
	}
	m.Mutex.Lock()
	m.owner = gid
	m.recursion++
}

func (m *RecursiveMutex) Unlock() {
	gid := utils.GoID()
	if m.owner != gid {
		panic(fmt.Errorf("wrong the owner(%d): %d", m.owner, gid))
	}
	m.recursion--
	if m.recursion != 0 {
		return
	}
	m.owner = 0
	m.Mutex.Unlock()

}
