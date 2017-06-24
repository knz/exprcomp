package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	testData := []struct {
		input string
		exp   string
	}{
		{"1.", "1"},
		{"1.2.+", "(1+2)"},
		{"12.23.+", "(12+23)"},
		{"a", "a()"},
		{"1.b", "b(1)"},
		{"", ""},    // error
		{"1.+", ""}, // error
	}
	for _, test := range testData {
		t.Run(test.input, func(t *testing.T) {
			node, err := parse(strings.NewReader(test.input))
			if err != nil {
				if test.exp != "" {
					t.Fatalf("%s: expected success, got %v", test.input, err)
				}
				return
			}
			nodeStr := node.String()
			if nodeStr != test.exp {
				t.Fatalf("%s: expected %s, got %s", test.input, test.exp, nodeStr)
			}
		})
	}
}
