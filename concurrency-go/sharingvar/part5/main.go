package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

const allLetters = "abcdefghijklmnopqrstuvwxyz"

type Task struct {
	url       string
	frequency []int
}

func main() {
	numWorkers := runtime.NumCPU()

	start := 1100
	end := 1180
	totalCount := (end - start) + 1

	tasks := make(chan string, totalCount)
	results := make(chan []int, totalCount)

	wg := &sync.WaitGroup{}

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(tasks, results, wg)
	}

	go func() {
		for i := start; i <= end; i++ {
			tasks <- fmt.Sprintf("https://rfc-editor.org/rfc/rfc%d.txt", i)
		}
		close(tasks)
	}()

	go func() {
		wg.Wait()
		close(results)
	}()

	//Aggregate results
	finalFrequency := make([]int, 26)
	for freq := range results {
		for i := 0; i < 26; i++ {
			finalFrequency[i] += freq[i]
		}
	}

	fmt.Println()
	for i, c := range allLetters {
		fmt.Printf("%c-%d\n", c, finalFrequency[i])
	}
	fmt.Println()
}

func worker(tasks chan string, results chan<- []int, wg *sync.WaitGroup) {
	defer wg.Done()
	for url := range tasks {
		frequency := make([]int, 26)
		err := countLetters(url, frequency)
		if err != nil {
			panic("worker failed to calculate count - " + err.Error())
		}
		results <- frequency
	}
}

func countLetters(url string, frequency []int) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("error while getting the data from %s: %v", url, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code %d for URL %s", resp.StatusCode, url)
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
	return nil
}
