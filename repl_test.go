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
			input:    "    hello    wOrld   ",
			expected: []string{"hello", "world"},
		},
		{
			input:    "  pikachu AnD raichu  not    BULBASAUR",
			expected: []string{"pikachu", "and", "raichu", "not", "bulbasaur"},
		},
		{
			input:    "Squirtle  is   S-TIER",
			expected: []string{"squirtle", "is", "s-tier"},
		},
	}

	for _, c := range cases {
		actual := cleanInput(c.input)
		if len(actual) != len(c.expected) {
			t.Errorf("cleanInput(%q):\nexpected: %d tokens\nactual: %d tokens\n\n", c.input, len(c.expected), len(actual))
			continue
		}
		for i := range actual {
			word := actual[i]
			expectedWord := c.expected[i]
			if word != expectedWord {
				t.Errorf("cleanInput(%q)[%d]:\nexpected: %q\nactual: %q\n\n", c.input, i, expectedWord, word)
			}
		}
	}
}
