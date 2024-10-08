package main

import (
	"fmt"
	"runtime"
)

func main() {
	// Total CPU core
	fmt.Println("Total no of CPU = ", runtime.NumCPU())
	// Used to set how many cores we want to use for go process
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
