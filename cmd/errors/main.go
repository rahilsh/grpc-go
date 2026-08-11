// Command errors runs the Go "errors & panic" examples: wrapped errors,
// panic/recover, and error wrapping with stack traces. Pick one with -demo:
//
//	go run ./cmd/errors -demo panic-recover
package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"custom-errors":  customErrors,
	"panic-recover":  panicRecover,
	"error-wrapping": errorWrapping,
}

// Config holds configuration (fields redacted for the example).
type Config struct{}

func readConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "can't open configuration file")
	}
	defer func() { _ = file.Close() }()
	return &Config{}, nil
}

func customErrors() {
	cfg, err := readConfig("/path/to/config.toml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		log.Printf("error with stack: %+v", err)
		return
	}
	fmt.Println(cfg)
}

func safeValue(vals []int, index int) (value int) {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("ERROR: %s\n", err)
		}
	}()
	return vals[index]
}

func panicRecover() {
	v := safeValue([]int{1, 2, 3}, 10) // out of range -> recovered
	fmt.Println(v)
}

func killServer(pidFile string) error {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		return errors.Wrap(err, "can't open pid file (is server running?)")
	}
	if err := os.Remove(pidFile); err != nil {
		log.Printf("warning: can't remove pid file - %s", err)
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return errors.Wrap(err, "bad process ID")
	}
	fmt.Printf("killing server with pid=%d\n", pid)
	return nil
}

func errorWrapping() {
	if err := killServer("server.pid"); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
	}
}
