package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	sequentialProcess()
	fmt.Println()
	parallelProcess()
}

func parallelProcess() {
	fmt.Println("Staring parallel process")
	start := time.Now()
	var wg sync.WaitGroup
	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go doWorkWithWG(i, &wg)
	}
	wg.Wait()
	fmt.Println("Total time for parallel process = ", int(time.Now().Sub(start).Seconds()))
}

func sequentialProcess() {
	fmt.Println("Staring sequential process")
	start := time.Now()
	for i := 1; i <= 5; i++ {
		doWork(i)
	}
	fmt.Println("Total time for sequential process = ", int(time.Now().Sub(start).Seconds()))
}

func doWorkWithWG(id int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(2 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}

func doWork(id int) {
	fmt.Printf("Work %d started at %s\n", id, time.Now().Format("15:04:05"))
	time.Sleep(2 * time.Second)
	fmt.Printf("Work %d finished at %s\n", id, time.Now().Format("15:04:05"))
}
