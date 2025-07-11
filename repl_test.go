package main

import (
	"testing"
)

func TestCleanInput(t *testing.T) {
	cases := []struct {
		input    string
		expected []string
	}{
		{
			input:    "    hello     world",
			expected: []string{"hello", "world"},
		},
		{
			input:    "PikaChu ChaRMANDER GEODUDE",
			expected: []string{"pikachu", "charmander", "geodude"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)

		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]

			if word != expectedWord {
				t.Errorf("cleanInput(%q) == %v, expected %v", c.input, actual, c.expected)
			}
		}
	}

}
