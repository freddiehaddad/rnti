package rnti

import (
	"log"
	"os"
	"strings"
	"sync"
)

func init() {
	log.SetOutput(os.Stderr)
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

// numeralValues is a mapping of the smallest Roman Numeral units to their
// integer representation.
var numeralValues = map[string]int{
	"I": 1,
	"V": 5,
	"X": 10,
	"L": 50,
	"C": 100,
	"D": 500,
	"M": 1000,

	"IV": 4,
	"IX": 9,
	"XL": 40,
	"XC": 90,
	"CD": 400,
	"CM": 900,
}

// readSymbols parses values from numeral and writes them to the numerals
// channel.  The numerals channel is closed when this function returns.
func readSymbols(numerals chan string, numeral string) {
	defer close(numerals)

	for i := 0; i < len(numeral); i++ {
		var sb strings.Builder

		ch := numeral[i]
		sb.WriteByte(ch)

		peek := byte(0)
		if i+1 < len(numeral) {
			peek = numeral[i+1]
		}

		switch ch {
		case 'I':
			switch peek {
			case 'V', 'X':
				sb.WriteByte(peek)
				i++
			}
		case 'X':
			switch peek {
			case 'L', 'C':
				sb.WriteByte(peek)
				i++
			}
		case 'C':
			switch peek {
			case 'D', 'M':
				sb.WriteByte(peek)
				i++
			}
		}

		log.Printf("readSymbols: symbol=%s", sb.String())
		numerals <- sb.String()
	}
}

// readValues reads Roman Numeral values from the numerals channel and looks up
// the numeric reprentation from the numeralValues map.  The value is written
// the values channel.  Upon returning (when the numerals channel is closed by
// readSymbols), the values channel is closed.
func readValues(numerals chan string, values chan int) {
	var wg sync.WaitGroup
	defer close(values)

	for numeral := range numerals {
		wg.Add(1)
		go func(n string) {
			log.Printf("readValues: processing numeral=%q", n)
			value := numeralValues[n]
			values <- value
			wg.Done()
		}(numeral)
	}

	wg.Wait()
}

// addValues reads integers from the values channel and sums them until the
// channel is closed.  The sum is written to the value channel which is closed
// upon returning.
func addValues(values chan int, value chan int) {
	defer close(value)

	sum := 0
	for v := range values {
		log.Printf("addValues: adding value=%d to sum=%d", v, sum)
		sum += v
	}

	value <- sum
}

// Convert reads the Roman Numeral string s and converts to an integer.  The
// function handles the conversion using multiple threads with the following
// tasks:
//
//   - Read lexical values from s
//   - Look up numeric value from the lexeme
//   - Sum the numeric values
func Convert(s string) int {
	numerals := make(chan string)
	values := make(chan int)
	value := make(chan int)

	go readSymbols(numerals, s)
	go readValues(numerals, values)
	go addValues(values, value)

	return <-value
}
