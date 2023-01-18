package reports

import (
	"log"
	"testing"
)

func Test_getCountStringFromOptimization(t *testing.T) {
	k := map[string][]string{
		"a": {"1", "2", "4", "5"},
		"b": {"1", "2", "3", "4", "5"},
		"c": {"1", "5"},
	}

	log.Print(getCountStringFromOptimization(k))
	t.Fail()
}

func TestNearest10(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected int
	}{
		{"less than 10", 5, 1},
		{"greater than or equal to 10 but less than 100", 22, 10},
		{"greater than or equal to 100 but less than 1000", 314, 100},
		{"greater than or equal to 1000 but less than 10000", 7234, 1000},
		{"greater than or equal to 10000", 12345, 10000},
		{"zero", 0, 1},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := nearest10(tc.input)
			if result != tc.expected {
				t.Errorf("For input %d, expected %d but got %d", tc.input, tc.expected, result)
			}
		})
	}
}
