package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var (
	url              string
	requestsQuantity int
	concurrency      int
)

var rootCmd = &cobra.Command{
	Use:   "st",
	Short: "Load stress test tool",
	Run: func(cmd *cobra.Command, args []string) {
		run(url, requestsQuantity, concurrency)
	},
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVar(&url, "url", "", "URL to be tested")
	rootCmd.Flags().IntVar(&requestsQuantity, "requests", 100, "number of requests to be sent")
	rootCmd.Flags().IntVar(&concurrency, "concurrency", 10, "number of concurrent requests")
	rootCmd.MarkFlagRequired("url")
}

func run(url string, requestsQuantity int, concurrency int) {
	var wg sync.WaitGroup
	requestChan := make(chan struct{}, concurrency)
	results := make(chan int, requestsQuantity)

	startTime := time.Now()
	fmt.Println("Benchmarking...")
	for i := 0; i < requestsQuantity; i++ {
		wg.Add(1)
		requestChan <- struct{}{}
		go func() {
			defer wg.Done()

			statusCode := makeRequest(url)
			results <- statusCode
			<-requestChan
		}()
	}

	wg.Wait()
	close(results)

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	statuses := make(map[int]int)
	for statusCode := range results {
		statuses[statusCode]++
	}

	fmt.Printf("Duration: %v\n", duration)
	fmt.Printf("Total of requests: %d\n", requestsQuantity)
	fmt.Printf("Total of http 200 responses: %d\n", statuses[200])
	for status, count := range statuses {
		if status != 200 {
			fmt.Printf("Total of http %d responses: %d\n", status, count)
		}
	}
}

func makeRequest(url string) int {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		return 0
	}

	resp, err := client.Do(req)
	if err != nil {
		return resp.StatusCode
	}
	defer resp.Body.Close()

	//just reading the body because if delete this line we couldn't see the results in an accurate way.
	io.ReadAll(resp.Body)

	return resp.StatusCode
}
