package main

import (
	"fmt"
	"runtime"
)

func main() {
	fmt.Println("Total no of CPU = ", runtime.NumCPU())
	fmt.Println("Go Max Process = ", runtime.GOMAXPROCS(0))
	contextSwitch()
}

func sayHello() {
	fmt.Println("Hello Go !!!")
}

func contextSwitch() {
	go sayHello()
	// Yields the context to sayHello
	runtime.Gosched()
	fmt.Println("Finished")
}
