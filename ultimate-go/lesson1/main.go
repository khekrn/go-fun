package main

import (
	"fmt"
	"unsafe"
)

func types() {
	var a int
	var b string
	var c float64
	var d bool

	var f uint = 1

	fmt.Printf("var a int \t %T [%v]\n", a, a)
	fmt.Printf("var b string \t %T [%v]\n", b, b)
	fmt.Printf("var c float64 \t %T [%v]\n", c, c)
	fmt.Printf("var d bool \t %T [%v]\n", d, d)
	fmt.Println("Int = ", f)
}

func padding() {
	// What will be size of this
	type fun struct {
		age       int16
		isMale    bool
		expenses  float32
		isMarried bool
		sal       float64
	}

	f := fun{}

	fmt.Printf("Size of fun = [%d]\n", unsafe.Sizeof(f))
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.age, &f.age)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.isMale, &f.isMale)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.expenses, &f.expenses)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.isMarried, &f.isMarried)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.sal, &f.sal)

	// A bool is 1 byte, int16 is 2 bytes, and float32 is 4 bytes.
	// Add that all together and I get 7 bytes. However, the actual answer
	// is 8 bytes. Why, because there is a padding byte sitting between
	// the flag and counter fields for the reason of alignment.
	type example struct {
		flag    bool
		counter int16
		pi      float32
	}

	fmt.Printf("Size of example = [%d]\n", unsafe.Sizeof(example{}))

	// What will be the size of this struct ?
	type example2 struct {
		flag    bool
		counter int16
		flag2   bool
		pi      float32
	}

	// Why 12 Byte ?
	fmt.Printf("Size of example2 = [%d]\n", unsafe.Sizeof(example2{}))
	// type example2 struct {
	// 	flag    bool	// 0xc000100020 <- Starting Address
	// 			byte	// 0xc000100021 <- 1 byte padding
	// 	counter int16	// 0xc000100022 <- 2 byte alignment
	// 	flag2   bool	// 0xc000100024 <- 1 byte alignment
	// 			byte	// 0xc000100025 <- 1 byte padding
	// 			byte	// 0xc000100026 <- 1 byte padding
	// 			byte	// 0xc000100027 <- 1 byte padding
	// 	pi      float32  // 0xc000100028 <- 4 byte alignment
	// 	}

	// This is how the alignment and padding play out if a value of type example2 starts at address
	// 0xc000100020. The flag field represents the starting address and is only 1 byte in size.
	// Since the counter field requires 2 bytes of allocation, it must be placed in memory on a 2-byte
	// alignment, meaning it needs to fall on an address that is a multiple of 2. This requires the
	// counter field to start at address 0xc000100022. This creates a 1-byte gap between the flag and
	// counter fields.

	// The flag2 field is a bool and can fall at the next address 0xc000100024. The final field is pi
	// and requires 4 bytes of allocation so it needs to fall on a 4-byte alignment. The next address
	// for a 4 byte value is at 0xc000100028. That means 3 more padding bytes are needed to maintain a
	// proper alignment. This results in a value of type example2 requiring 12 bytes of total memory
	// allocation.

	// One can save the 4 bytes just following simple pattern when defining struct fields
	// Allocate the fields from highest to lowest
	type example3 struct {
		pi      float32
		counter int16
		flag    bool
		flag2   bool
	}

	fmt.Printf("Size of example3 = [%d]\n", unsafe.Sizeof(example3{}))
	// Now it's 8 byte but how ?

	// The largest field in a struct represents the alignment boundary for the entire struct.
	// In this case, the largest field is 4 bytes so the starting address for this struct value must
	// be a multiple of 4. I can see the address 0xc000100020 is a multiple of 4.

	// If I need to minimize the amount of padding bytes, I must lay out the fields from highest
	// allocation to lowest allocation. This will push any necessary padding bytes down to the bottom
	// of the struct and reduce the total number of padding bytes necessary.

	// type example3 struct {
	// 	pi      float32 // 0xc000100020 <- Starting Address
	// 	counter int16   // 0xc000100024 <- 2 byte alignment
	// 	flag    bool    // 0xc000100026 <- 1 byte alignment
	// 	flag2   bool    // 0xc000100027 <- 1 byte alignment
	// }
}

func main() {
	padding()
	types()
	fmt.Println("Padding\n")
	homework()
}

func homework() {
	type s1 struct {
		flag   bool
		age    int32
		salary float64
		gender int16
	}

	f := &s1{false, 12, 32423476436783467853476.546, 1}

	fmt.Printf("Value[%v]\tAddress[%p]\n", f.flag, &f.flag)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.age, &f.age)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.salary, &f.salary)
	fmt.Printf("Value[%v]\tAddress[%p]\n", f.gender, &f.gender)
	fmt.Printf("Size of S1 = [%d]\n", unsafe.Sizeof(f))
}
