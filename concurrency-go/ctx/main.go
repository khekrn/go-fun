package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	generator := func(dataItem string, stream chan any) {
		for {
			select {
			case <-ctx.Done():
				return
			case stream <- dataItem:
			}
		}
	}
	infiniteApples := make(chan any)
	go generator("Apple", infiniteApples)

	infiniteOranges := make(chan any)
	go generator("Orange", infiniteOranges)

	infinitePeaches := make(chan any)
	go generator("Peach", infinitePeaches)

	wg.Add(1)
	go spawn(ctx, &wg, infiniteApples)

	wg.Add(1)
	go genericFunc(ctx, &wg, infiniteOranges)

	wg.Add(1)
	go genericFunc(ctx, &wg, infinitePeaches)

	wg.Wait()
}

func spawn(ctx context.Context, parentWg *sync.WaitGroup, infiniteApples chan any) {
	defer parentWg.Done()
	var wg sync.WaitGroup
	doWork := func(ctx context.Context) {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-infiniteApples:
				if !ok {
					fmt.Println("channel closed")
					return
				}
				fmt.Println(d)
			}
		}
	}
	newCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go doWork(newCtx)
	}
	wg.Wait()
}

func genericFunc(ctx context.Context, wg *sync.WaitGroup, stream chan any) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-stream:
			if !ok {
				fmt.Println("channel closed")
				return
			}
			fmt.Println(d)
		}
	}
}
