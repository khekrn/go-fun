package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

// Sequential code
func main() {
	mutex := sync.Mutex{}
	var frequency = make([]int, 26)
	for i := 1000; i <= 1030; i++ {
		url := fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		go countLetters(url, frequency, &mutex)
	}

	time.Sleep(20 * time.Second)
	fmt.Println()
	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, frequency[i])
	}
	fmt.Println()
}

func countLetters(url string, frequency []int, mutex *sync.Mutex) {
	resp, err := http.Get(url)
	if err != nil {
		panic("error while getting the data from " + url + " with an error = " + err.Error())
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		panic("cannot read from the given url " + url)
	}

	mutex.Lock()
	defer mutex.Unlock()
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
