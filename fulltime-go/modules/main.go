package main

import "fmt"

func main() {
	fmt.Println("The Number is : ", getNumbers())
}

// Notes
// If you try to run `go run main.go` you will get the following error
// ./main.go:6:34: undefined: getNumbers
// This is because go run only takes one file where as here
// we have two files, main.go and numbers.go
// In order to fix this you can create the module using
// `go mod init {name}` and `go build -o {name}` which will build
// the entire go files
