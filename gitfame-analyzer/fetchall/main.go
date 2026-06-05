//go:build !solution

package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

func main() {
	urls := os.Args[1:]
	if len(urls) == 0 {
		fmt.Fprintln(os.Stderr, "usage: fetchall <url1> <url2> ...")
		os.Exit(1)
	}
	start := time.Now()
	client := &http.Client{Timeout: 15 * time.Second}
	results := make(chan string, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))

	for _, u := range urls {
		url := u
		go func() {
			defer wg.Done()
			fetch(url, client, results)
		}()
	}
	go func() {
		wg.Wait()
		close(results)
	}()
	for line := range results {
		fmt.Println(line)
	}
	fmt.Printf("%.2fs elapsed\n", time.Since(start).Seconds())
}
func fetch(url string, client *http.Client, out chan<- string) {
	t0 := time.Now()
	resp, err := client.Get(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get %s: %v\n", url, err)
		out <- fmt.Sprintf("%.2fs    %7d  %s", time.Since(t0).Seconds(), 0, url)
		return
	}
	defer resp.Body.Close()

	n, copyErr := io.Copy(io.Discard, resp.Body)
	if copyErr != nil {
		fmt.Fprintf(os.Stderr, "Read %s: %v\n", url, copyErr)
	}
	out <- fmt.Sprintf("%.2fs    %7d  %s", time.Since(t0).Seconds(), n, url)
}
