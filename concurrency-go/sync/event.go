package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mx   = sync.Mutex{}
	cond = sync.NewCond(&mx)

	events = make(map[string]bool)
)

func waiter(name, event string) {
	mx.Lock()
	for !events[event] {
		cond.Wait()
	}
	fmt.Println("Name = ", name, " and Event = ", event)
	mx.Unlock()
}

func signalEvent(event string) {
	time.Sleep(2 * time.Second)
	mx.Lock()
	events[event] = true
	cond.Broadcast()
	mx.Unlock()
}

func main() {
	go waiter("Goroutine 1", "EventA")
	go waiter("Goroutine 2", "EventA")
	go waiter("Goroutine 3", "EventB")
	go waiter("Goroutine 4", "EventB")

	go signalEvent("EventA")
	go signalEvent("EventB")

	time.Sleep(5 * time.Second)
}
