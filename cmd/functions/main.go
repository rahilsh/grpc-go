// Command functions runs the Go "functions" examples: definitions, parameter
// passing, error returns, defer, and an HTTP helper. Pick one with -demo, e.g.:
//
//	go run ./cmd/functions -demo defer
package main

import (
	"fmt"
	"math"
	"net/http"

	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"functions":           functions,
	"function-parameters": functionParameters,
	"returning-errors":    returningErrors,
	"defer":               deferOrder,
	"content-type":        contentTypeDemo,
}

func add(a, b int) int           { return a + b }
func divmod(a, b int) (int, int) { return a / b, a % b }

func functions() {
	fmt.Println(add(1, 2))
	div, mod := divmod(7, 2)
	fmt.Printf("div=%d, mod=%d\n", div, mod)
}

func doubleAt(values []int, i int) { values[i] *= 2 }
func double(n int)                 { n *= 2; fmt.Println("inside double (local copy):", n) }
func doublePtr(n *int)             { *n *= 2 }

func functionParameters() {
	values := []int{1, 2, 3, 4}
	doubleAt(values, 2)
	fmt.Println(values)

	val := 10
	double(val)
	fmt.Println(val)
	doublePtr(&val)
	fmt.Println(val)
}

func sqrt(n float64) (float64, error) {
	if n < 0 {
		return 0.0, fmt.Errorf("sqrt of negative value (%f)", n)
	}
	return math.Sqrt(n), nil
}

func returningErrors() {
	for _, n := range []float64{2.0, -2.0} {
		if s, err := sqrt(n); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		} else {
			fmt.Println(s)
		}
	}
}

func cleanup(name string) { fmt.Printf("Cleaning up %s\n", name) }

func deferOrder() {
	defer cleanup("A")
	defer cleanup("B")
	fmt.Println("worker")
}

// contentType returns the Content-Type header of an HTTP GET to url.
func contentType(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	ctype := resp.Header.Get("Content-Type")
	if ctype == "" {
		return "", fmt.Errorf("can't find Content-Type header")
	}
	return ctype, nil
}

func contentTypeDemo() {
	if ctype, err := contentType("https://linkedin.com"); err != nil {
		fmt.Printf("ERROR: %s\n", err)
	} else {
		fmt.Println(ctype)
	}
}
