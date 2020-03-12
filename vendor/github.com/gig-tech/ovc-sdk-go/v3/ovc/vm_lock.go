package ovc

import (
	"sync"
)

var (
	lock         = &sync.Mutex{}
	locks        = make(map[int]*sync.Mutex)
	lockCounters = make(map[int]int)
)

// GetLock returns when its safe to execute a synchronized action towards a certain vm
func GetLock(vmID int) {
	lock.Lock()
	vmLock := locks[vmID]
	if vmLock == nil {
		vmLock = &sync.Mutex{}
		locks[vmID] = vmLock
		lockCounters[vmID] = 0
	}
	lockCounters[vmID]++
	lock.Unlock()
	vmLock.Lock()
}

// ReleaseLock free's access to execute a certain action towards a certain vm
func ReleaseLock(vmID int) {
	defer lock.Unlock()
	lock.Lock()
	vmLock := locks[vmID]
	vmLock.Unlock()
	lockCounters[vmID]--
	if lockCounters[vmID] == 0 {
		delete(locks, vmID)
		delete(lockCounters, vmID)
	}
}
