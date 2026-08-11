// Command basics runs the Go language "basics" examples: variables, control
// flow, strings, slices, and maps. Pick one with -demo, e.g.:
//
//	go run ./cmd/basics -demo fizzbuzz
package main

import (
	"fmt"
	"strings"

	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"hello-world":        helloWorld,
	"mean":               mean,
	"if":                 ifStatement,
	"switch":             switchStatement,
	"loops":              loops,
	"fizzbuzz":           fizzBuzz,
	"strings":            stringValues,
	"sprintf":            sprintf,
	"even-ended-numbers": evenEndedNumbers,
	"slices":             slices,
	"slice-max":          sliceMax,
	"maps":               maps,
	"word-count":         wordCount,
}

func helloWorld() {
	fmt.Println("Welcome Gophers ☺")
}

func mean() {
	x := 1.0
	y := 2.0
	fmt.Printf("x=%v, type of %T\n", x, x)
	fmt.Printf("y=%v, type of %T\n", y, y)
	m := (x + y) / 2.0
	fmt.Printf("result: %v, type of %T\n", m, m)
}

func ifStatement() {
	x := 10
	if x > 5 {
		fmt.Println("x is big")
	}
	if x > 100 {
		fmt.Println("x is very big")
	} else {
		fmt.Println("x is not that big")
	}
	if x > 5 && x < 15 {
		fmt.Println("x is just right")
	}
	if x < 20 || x > 30 {
		fmt.Println("x is out of range")
	}
	a, b := 11.0, 20.0
	if frac := a / b; frac > 0.5 {
		fmt.Println("a is more than half of b")
	}
}

func switchStatement() {
	x := 2
	switch x {
	case 1:
		fmt.Println("one")
	case 2:
		fmt.Println("two")
	case 3:
		fmt.Println("three")
	default:
		fmt.Println("many")
	}
	switch {
	case x > 100:
		fmt.Println("x is very big")
	case x > 10:
		fmt.Println("x is big")
	default:
		fmt.Println("x is small")
	}
}

func loops() {
	for i := 0; i < 3; i++ {
		fmt.Println(i)
	}
	fmt.Println("----")
	for i := 0; i < 3; i++ {
		if i > 1 {
			break
		}
		fmt.Println(i)
	}
	fmt.Println("----")
	for i := 0; i < 3; i++ {
		if i < 1 {
			continue
		}
		fmt.Println(i)
	}
	fmt.Println("----")
	a := 0
	for a < 3 {
		fmt.Println(a)
		a++
	}
}

func fizzBuzz() {
	for i := 1; i <= 20; i++ {
		switch {
		case i%3 == 0 && i%5 == 0:
			fmt.Println("fizz buzz")
		case i%3 == 0:
			fmt.Println("fizz")
		case i%5 == 0:
			fmt.Println("buzz")
		default:
			fmt.Println(i)
		}
	}
}

func stringValues() {
	book := "The colour of magic"
	fmt.Println(book)
	fmt.Println(len(book))
	fmt.Printf("book[0] = %v (type %T)\n", book[0], book[0]) // byte
	fmt.Println(book[4:11])
	fmt.Println(book[4:])
	fmt.Println(book[:4])
	fmt.Println("t" + book[1:])
	fmt.Println("It was ½ price!") // unicode
	poem := "\n\tThe road goes ever on\n\t..."
	fmt.Println(poem)
}

func sprintf() {
	n := 42
	s := fmt.Sprintf("%d", n)
	fmt.Printf("s = %q (type %T)\n", s, s)
}

// evenEndedNumbers counts products of two 4-digit numbers whose first and last
// digit are the same.
func evenEndedNumbers() {
	count := 0
	for a := 1000; a <= 9999; a++ {
		for b := a; b <= 9999; b++ {
			s := fmt.Sprintf("%d", a*b)
			if s[0] == s[len(s)-1] {
				count++
			}
		}
	}
	fmt.Println(count)
}

func slices() {
	loons := []string{"bugs", "daffy", "taz"}
	fmt.Printf("loons = %v (type %T)\n", loons, loons)
	fmt.Println(len(loons))
	fmt.Println(loons[1])
	fmt.Println(loons[1:])
	for i, name := range loons {
		fmt.Printf("%s at %d\n", name, i)
	}
	loons = append(loons, "elmer")
	fmt.Println(loons)
}

func sliceMax() {
	nums := []int{16, 8, 42, 4, 23, 15}
	largest := nums[0]
	for _, value := range nums[1:] {
		if value > largest {
			largest = value
		}
	}
	fmt.Println(largest)
}

func maps() {
	stocks := map[string]float64{
		"AMZN": 1699.8,
		"GOOG": 1129.19,
		"MSFT": 98.61,
	}
	fmt.Println(len(stocks))
	fmt.Println(stocks["MSFT"])
	if value, ok := stocks["TSLA"]; !ok {
		fmt.Println("TSLA not found")
	} else {
		fmt.Println(value)
	}
	stocks["TSLA"] = 322.12
	delete(stocks, "AMZN")
	for key, value := range stocks {
		fmt.Printf("%s -> %.2f\n", key, value)
	}
}

func wordCount() {
	text := "Needles and pins\nNeedles and pins\nSew me a sail"
	counts := map[string]int{}
	for _, word := range strings.Fields(text) {
		counts[strings.ToLower(word)]++
	}
	fmt.Println(counts)
}
