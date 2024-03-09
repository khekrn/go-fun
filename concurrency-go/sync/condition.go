package main

import (
	"fmt"
	"sync"
	"time"
)

type Buffer struct {
	sync.Mutex
	cond    *sync.Cond
	buffer  []int
	maxSize int
}

func NewBuffer(size int) *Buffer {
	buffer := Buffer{
		buffer:  make([]int, 0, size),
		maxSize: size,
	}
	buffer.cond = sync.NewCond(&buffer.Mutex)
	return &buffer
}

func (b *Buffer) Produce(item int) {
	b.Lock()
	defer b.Unlock()

	for len(b.buffer) == b.maxSize {
		b.cond.Wait()
	}

	b.buffer = append(b.buffer, item)
	fmt.Println("Item Produced = ", item)
	b.cond.Signal()
}

func (b *Buffer) Consume() int {
	b.Lock()
	defer b.Unlock()
	for len(b.buffer) == 0 {
		b.cond.Wait()
	}
	item := b.buffer[0]
	b.buffer = b.buffer[1:]

	b.cond.Signal()
	return item
}

func main() {
	buffer := NewBuffer(5)
	fmt.Println("Starting the goroutines")
	go func() {
		for i := range 10 {
			buffer.Produce(i)
			time.Sleep(500 * time.Millisecond)
		}
	}()

	go func() {
		for i := range 10 {
			result := buffer.Consume()
			fmt.Println("For Iteration = ", i, " Item Consumed = ", result)
			time.Sleep(1000 * time.Millisecond)
		}
	}()

	time.Sleep(15 * time.Second)
	fmt.Println("Program Completed")
}
