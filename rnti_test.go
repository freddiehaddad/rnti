package rnti

import (
	"testing"
)

func TestRomanNumeralMap(t *testing.T) {
	tests := []struct {
		numeral string
		value   int
	}{
		{"I", 1},
		{"V", 5},
		{"X", 10},
		{"L", 50},
		{"C", 100},
		{"D", 500},
		{"M", 1000},

		{"IV", 4},
		{"IX", 9},
		{"XL", 40},
		{"XC", 90},
		{"CD", 400},
		{"CM", 900},
	}

	for i, test := range tests {
		value, ok := numeralValues[test.numeral]
		if !ok {
			t.Errorf("Test[%d]: Roman numeral %q missing from map", i, test.numeral)
		}

		if test.value != value {
			t.Errorf("Test[%d]: Value for roman numeral %q wrong. Expected=%d Got=%d",
				i, test.numeral, test.value, value)
		}
	}
}

func TestReadSymbols(t *testing.T) {
	input := "MDCLXVIIVIXXLXCCDCM"
	expected := []string{"M", "D", "C", "L", "X", "V", "I", "IV", "IX", "XL", "XC", "CD", "CM"}

	numerals := make(chan string)

	go readSymbols(numerals, input)

	for i, e := range expected {
		numeral := <-numerals
		if e != numeral {
			t.Errorf("Test[%d]: Numeral wrong. Expected=%q Got=%q", i, e, numeral)
		}
	}
}

func TestReadValues(t *testing.T) {
	input := []string{"M", "D", "C", "L", "X", "V", "I", "IV", "IX", "XL", "XC", "CD", "CM"}
	expected := []int{1000, 500, 100, 50, 10, 5, 1, 4, 9, 40, 90, 400, 900}

	numerals := make(chan string)
	values := make(chan int)

	go func() {
		for _, numeral := range input {
			numerals <- numeral
		}
	}()

	go readValues(numerals, values)

	// Values arrive asynchronously. So we cannot guarantee the ordering.
	// Thus, it is necessary to search the expected array for the newly
	// arrived value and delete it.
	for len(expected) > 0 {
		value := <-values
		// Find the value in the expected array.
		index := func() int {
			for i := 0; i < len(expected); i++ {
				if expected[i] == value {
					return i
				}
			}
			return -1
		}()

		if index == -1 {
			t.Errorf("Expected value=%d not found", value)
			t.Fail()
		} else {
			expected[index] = expected[len(expected)-1]
			expected = expected[:len(expected)-1]
		}
	}
}

func TestAddValues(t *testing.T) {
	input := []int{1000, 500, 100, 50, 10, 5, 1, 4, 9, 40, 90, 400, 900}
	expected := 1000 + 500 + 100 + 50 + 10 + 5 + 1 + 4 + 9 + 40 + 90 + 400 + 900

	values := make(chan int)
	value := make(chan int)

	go func() {
		defer close(values)
		for _, val := range input {
			values <- val
		}
	}()

	go addValues(values, value)

	result := <-value
	if result != expected {
		t.Errorf("Failed: Value wrong. Expected=%d Got=%d", expected, result)
	}
}

func TestConvert(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"I", 1},
		{"II", 2},
		{"II", 2},
		{"III", 3},
		{"IIII", 4},
		{"IIIV", 6},
		{"V", 5},
		{"X", 10},
		{"L", 50},
		{"C", 100},
		{"D", 500},
		{"M", 1000},
		{"IV", 4},
		{"IX", 9},
		{"XL", 40},
		{"XC", 90},
		{"CD", 400},
		{"CM", 900},
		{"XIII", 13},
		{"LVII", 57},
		{"MCMXCIV", 1994},
	}

	for i, test := range tests {
		result := Convert(test.input)
		if result != test.expected {
			t.Errorf("Test[%d]: Result wrong. Expected=%d Got=%d", i, test.expected, result)
		}
	}
}
