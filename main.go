package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"sync"
	"time"
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
	fmt.Printf("Total of http 200 responses: %d\n", statuses[200])
}

func makeRequest(url string) int {
	req, err := http.Get(url)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return 0
	}
	_ = req.Body.Close()
	return req.StatusCode
}
