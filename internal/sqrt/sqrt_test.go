package sqrt

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"testing"
)

func almostEqual(v1, v2 float64) bool {
	return Abs(v1-v2) <= 0.001
}

// TestSimple is a single, hand-written assertion.
func TestSimple(t *testing.T) {
	val, err := Sqrt(2)
	if err != nil {
		t.Fatalf("error in calculation - %s", err)
	}
	if !almostEqual(val, 1.414214) {
		t.Fatalf("bad value - %f", val)
	}
}

// TestMany is a table-driven test with inline cases.
func TestMany(t *testing.T) {
	testCases := []struct {
		value    float64
		expected float64
	}{
		{0.0, 0.0},
		{2.0, 1.414214},
		{9.0, 3.0},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("%f", tc.value), func(t *testing.T) {
			out, err := Sqrt(tc.value)
			if err != nil {
				t.Fatal(err)
			}
			if !almostEqual(out, tc.expected) {
				t.Fatalf("%f != %f", out, tc.expected)
			}
		})
	}
}

// TestFromCSV is a table-driven test whose cases are loaded from sqrt_cases.csv.
func TestFromCSV(t *testing.T) {
	file, err := os.Open("sqrt_cases.csv")
	if err != nil {
		t.Fatalf("can't open cases file - %s", err)
	}
	defer func() { _ = file.Close() }()

	rdr := csv.NewReader(file)
	for {
		record, err := rdr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			t.Fatalf("error reading cases file - %s", err)
		}

		val, err := strconv.ParseFloat(record[0], 64)
		if err != nil {
			t.Fatalf("bad value - %s", record[0])
		}
		expected, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			t.Fatalf("bad value - %s", record[1])
		}

		t.Run(fmt.Sprintf("%f", val), func(t *testing.T) {
			out, err := Sqrt(val)
			if err != nil {
				t.Fatal(err)
			}
			if !almostEqual(out, expected) {
				t.Fatalf("%f != %f", out, expected)
			}
		})
	}
}

// BenchmarkSqrt measures Sqrt performance.
func BenchmarkSqrt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		if _, err := Sqrt(float64(i)); err != nil {
			b.Fatal(err)
		}
	}
}
