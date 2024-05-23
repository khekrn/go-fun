package main

import (
	"fmt"
	"sync"
	"time"
)

type ReadWriteMutex struct {
	*sync.Cond
	readCounter    int
	writersWaiting int
	writerActive   bool
}

func NewReadWriteMutex() *ReadWriteMutex {
	return &ReadWriteMutex{Cond: sync.NewCond(&sync.Mutex{})}
}

func (rwm *ReadWriteMutex) ReadLock() {
	rwm.L.Lock()
	for rwm.writerActive || rwm.writersWaiting > 0 {
		rwm.Wait()
	}
	rwm.readCounter++
	rwm.L.Unlock()
}

func (rwm *ReadWriteMutex) WriteLock() {
	rwm.L.Lock()
	rwm.writersWaiting++
	for rwm.readCounter > 0 || rwm.writerActive {
		rwm.Wait()
	}
	rwm.writersWaiting--
	rwm.writerActive = true
	rwm.L.Unlock()
}

func (rwm *ReadWriteMutex) ReadUnLock() {
	rwm.L.Lock()
	rwm.readCounter--
	if rwm.readCounter == 0 {
		rwm.Broadcast()
	}
	rwm.L.Unlock()
}
func (rwm *ReadWriteMutex) WriteUnLock() {
	rwm.L.Lock()
	rwm.writerActive = false
	rwm.Broadcast()
	rwm.L.Unlock()
}

func main() {
	rwMutex := NewReadWriteMutex()
	for i := 0; i < 2; i++ {
		go func() {
			for {
				rwMutex.ReadLock()
				time.Sleep(1 * time.Second)
				fmt.Println("Read done")
				rwMutex.ReadUnLock()
			}
		}()
	}
	time.Sleep(1 * time.Second)
	rwMutex.WriteLock()
	fmt.Println("Write finished")
	rwMutex.WriteUnLock()
}
