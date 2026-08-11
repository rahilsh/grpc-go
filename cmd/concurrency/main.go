// Command concurrency runs the Go concurrency examples: goroutines, channels,
// select, and a concurrent md5 checker. Pick one with -demo, e.g.:
//
//	go run ./cmd/concurrency -demo select
package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"goroutines":           goroutines,
	"channels":             channels,
	"channel-content-type": channelContentType,
	"select":               selectStmt,
	"md5-concurrent":       md5Concurrent,
}

func fetchType(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("error: %s\n", err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	fmt.Printf("%s -> %s\n", url, resp.Header.Get("content-type"))
}

func goroutines() {
	urls := []string{"https://golang.org", "https://api.github.com", "https://httpbin.org/xml"}
	var wg sync.WaitGroup
	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			fetchType(url)
		}(url)
	}
	wg.Wait()
}

func channels() {
	ch := make(chan int)
	go func() { ch <- 353 }()
	fmt.Printf("got %d\n", <-ch)

	fmt.Println("-----")
	go func() {
		for i := 0; i < 3; i++ {
			ch <- i
			time.Sleep(100 * time.Millisecond)
		}
		close(ch)
	}()
	for i := range ch {
		fmt.Printf("received %d\n", i)
	}
}

func fetchTypeChan(url string, out chan string) {
	resp, err := http.Get(url)
	if err != nil {
		out <- fmt.Sprintf("%s -> error: %s", url, err)
		return
	}
	defer func() { _ = resp.Body.Close() }()
	out <- fmt.Sprintf("%s -> %s", url, resp.Header.Get("content-type"))
}

func channelContentType() {
	urls := []string{"https://golang.org", "https://api.github.com", "https://httpbin.org/xml"}
	ch := make(chan string)
	for _, url := range urls {
		go fetchTypeChan(url, ch)
	}
	for range urls {
		fmt.Println(<-ch)
	}
}

func selectStmt() {
	ch1, ch2 := make(chan int), make(chan int)
	go func() { ch1 <- 42 }()
	select {
	case val := <-ch1:
		fmt.Printf("got %d from ch1\n", val)
	case val := <-ch2:
		fmt.Printf("got %d from ch2\n", val)
	}

	fmt.Println("----")
	out := make(chan float64)
	go func() {
		time.Sleep(100 * time.Millisecond)
		out <- 3.14
	}()
	select {
	case val := <-out:
		fmt.Printf("got %f\n", val)
	case <-time.After(20 * time.Millisecond):
		fmt.Println("timeout")
	}
}

func parseSignaturesFile(path string) (map[string]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	sigs := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for lnum := 1; scanner.Scan(); lnum++ {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			return nil, fmt.Errorf("%s:%d bad line", path, lnum)
		}
		sigs[fields[1]] = fields[0]
	}
	return sigs, scanner.Err()
}

func fileMD5(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func() { _ = file.Close() }()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

type result struct {
	path  string
	match bool
	err   error
}

func md5Worker(path, sig string, out chan *result) {
	r := &result{path: path}
	s, err := fileMD5(path)
	if err != nil {
		r.err = err
	} else {
		r.match = s == sig
	}
	out <- r
}

func md5Concurrent() {
	sigs, err := parseSignaturesFile("md5sum.txt")
	if err != nil {
		fmt.Printf("error: can't read signature file - %s\n", err)
		return
	}

	out := make(chan *result)
	for path, sig := range sigs {
		go md5Worker(path, sig, out)
	}
	for range sigs {
		switch r := <-out; {
		case r.err != nil:
			fmt.Printf("%s: error - %s\n", r.path, r.err)
		case !r.match:
			fmt.Printf("%s: signature mismatch\n", r.path)
		default:
			fmt.Printf("%s: ok\n", r.path)
		}
	}
}
