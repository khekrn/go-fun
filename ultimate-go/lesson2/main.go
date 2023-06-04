package main

import "fmt"

func arrays() {
	five := [5]string{"Kishore", "Aswin", "Rahul", "Varun", "Surya"}
	for i, v := range five {
		fmt.Printf("Value[%s]\tAddress[%p] IndexAddress[%p]\n", v, &v, &five[i])
	}
}

func main() {
	fmt.Println("Arrays are Contiguous")
	arrays()
}
