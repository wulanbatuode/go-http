package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitAndFilter(t *testing.T) {
	testCases := []struct {
		Str    string
		Sep    string
		Tokens []string
	}{
		{"abc: def", ":", []string{"abc", "def"}},
		{"abc : def", ":", []string{"abc", "def"}},
		{"//", "/", []string{}},
		{"/", "/", []string{}},
		{"/a/b//c", "/", []string{"a", "b", "c"}},
		{"a/b//c", "/", []string{"a", "b", "c"}},
	}
	for _, testCase := range testCases {
		assert.Equal(t, testCase.Tokens, SplitConvertFilter(testCase.Str, testCase.Sep,
			func(s string) string { return strings.Trim(s, " ") },
			func(s string) bool { return s != "" }))
	}

}
