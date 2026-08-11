// Package demo provides a tiny command-line dispatcher shared by the example
// command binaries. Each binary registers a set of named demo functions and
// Run executes the one selected with the -demo flag.
package demo

import (
	"flag"
	"fmt"
	"os"
	"sort"
)

// Run parses the -demo flag and invokes the matching function. When the flag is
// empty or names an unknown demo, it prints the available demo names.
func Run(demos map[string]func()) {
	name := flag.String("demo", "", "which demo to run")
	flag.Parse()

	if fn, ok := demos[*name]; ok {
		fn()
		return
	}

	if *name != "" {
		fmt.Fprintf(os.Stderr, "unknown demo: %q\n\n", *name)
	}
	fmt.Fprintln(os.Stderr, "available demos:")
	for _, n := range sortedKeys(demos) {
		fmt.Fprintf(os.Stderr, "  %s\n", n)
	}
	if *name != "" {
		os.Exit(1)
	}
}

func sortedKeys(demos map[string]func()) []string {
	names := make([]string, 0, len(demos))
	for n := range demos {
		names = append(names, n)
	}
	sort.Strings(names)
	return names
}
