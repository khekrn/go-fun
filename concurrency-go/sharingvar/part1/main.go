package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

// Sequential code
func main() {
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		countLetters(url, frequency)
	}

	fmt.Println()
	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, frequency[i])
	}
	fmt.Println()
}

func countLetters(url string, frequency []int) {
	resp, err := http.Get(url)
	if err != nil {
		panic("error while getting the data from " + url + " with an error = " + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("cannot read from the given url " + url)
	}

	body, _ := io.ReadAll(resp.Body)
	for _, v := range body {
		c := strings.ToLower(string(v))
		cIndex := strings.Index(allLetters, c)
		if cIndex >= 0 {
			frequency[cIndex] += 1
		}
	}
	fmt.Println("Completed : ", url)
}
