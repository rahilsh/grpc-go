// Command types runs the Go "types & interfaces" examples: structs, methods,
// constructors, embedding, interfaces, and io.Writer. Pick one with -demo, e.g.:
//
//	go run ./cmd/types -demo interfaces
package main

import (
	"fmt"
	"io"
	"log"
	"math"
	"os"

	"github.com/rahilsh/golang-lab/internal/demo"
)

func main() { demo.Run(demos) }

var demos = map[string]func(){
	"structs":          structs,
	"receivers":        receivers,
	"methods":          methods,
	"constructor":      constructor,
	"embedded-structs": embeddedStructs,
	"interfaces":       interfaces,
	"io-writer":        ioWriter,
}

// Trade is a trade in stocks.
type Trade struct {
	Symbol string
	Volume int
	Price  float64
	Buy    bool
}

// NewTrade creates a trade after validating the input.
func NewTrade(symbol string, volume int, price float64, buy bool) (*Trade, error) {
	switch {
	case symbol == "":
		return nil, fmt.Errorf("symbol can't be empty")
	case volume <= 0:
		return nil, fmt.Errorf("volume must be >= 0 (was %d)", volume)
	case price <= 0.0:
		return nil, fmt.Errorf("price must be >= 0 (was %f)", price)
	}
	return &Trade{Symbol: symbol, Volume: volume, Price: price, Buy: buy}, nil
}

// Value returns the trade value (negative for a buy).
func (t *Trade) Value() float64 {
	value := float64(t.Volume) * t.Price
	if t.Buy {
		value = -value
	}
	return value
}

func structs() {
	t1 := Trade{"MSFT", 10, 99.98, true}
	fmt.Println(t1)
	fmt.Printf("%+v\n", t1)
	fmt.Println(t1.Symbol)
	fmt.Printf("%+v\n", Trade{})
}

func methods() {
	t := Trade{Symbol: "MSFT", Volume: 10, Price: 99.98, Buy: true}
	fmt.Println(t.Value())
}

func constructor() {
	t, err := NewTrade("MSFT", 10, 99.98, true)
	if err != nil {
		fmt.Printf("error: can't create trade - %s\n", err)
		return
	}
	fmt.Println(t.Value())
}

// Point is a 2D point.
type Point struct {
	X int
	Y int
}

// Move translates the point.
func (p *Point) Move(dx, dy int) {
	p.X += dx
	p.Y += dy
}

func receivers() {
	p := &Point{1, 2}
	p.Move(2, 3)
	fmt.Printf("%+v\n", p)
}

// Square is a square built by embedding a Point as its center.
type Square struct {
	Center Point
	Length int
}

// NewSquare creates a square after validating the length.
func NewSquare(x, y, length int) (*Square, error) {
	if length <= 0 {
		return nil, fmt.Errorf("length must be > 0")
	}
	return &Square{Center: Point{x, y}, Length: length}, nil
}

// Move translates the square by moving its center.
func (s *Square) Move(dx, dy int) { s.Center.Move(dx, dy) }

// Area returns the square's area.
func (s *Square) Area() int { return s.Length * s.Length }

func embeddedStructs() {
	s, err := NewSquare(1, 1, 10)
	if err != nil {
		log.Fatalf("ERROR: can't create square")
	}
	s.Move(2, 3)
	fmt.Printf("%+v\n", s)
	fmt.Println(s.Area())
}

// Shape is anything with an area.
type Shape interface {
	Area() float64
}

// ShapeSquare is a square expressed purely by its side length.
type ShapeSquare struct {
	Length float64
}

// Area returns the area of the square.
func (s *ShapeSquare) Area() float64 { return s.Length * s.Length }

// Circle is a circle.
type Circle struct {
	Radius float64
}

// Area returns the area of the circle.
func (c *Circle) Area() float64 { return math.Pi * c.Radius * c.Radius }

func sumAreas(shapes []Shape) float64 {
	total := 0.0
	for _, shape := range shapes {
		total += shape.Area()
	}
	return total
}

func interfaces() {
	s := &ShapeSquare{20}
	fmt.Println(s.Area())
	c := &Circle{10}
	fmt.Println(c.Area())
	fmt.Println(sumAreas([]Shape{s, c}))
}

// Capper is an io.Writer that upper-cases everything written through it.
type Capper struct {
	wtr io.Writer
}

func (c *Capper) Write(p []byte) (int, error) {
	diff := byte('a' - 'A')
	out := make([]byte, len(p))
	for i, b := range p {
		if b >= 'a' && b <= 'z' {
			b -= diff
		}
		out[i] = b
	}
	return c.wtr.Write(out)
}

func ioWriter() {
	c := &Capper{os.Stdout}
	_, _ = fmt.Fprintln(c, "Hello there")
}
