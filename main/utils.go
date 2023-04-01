package main

import "strings"

func SplitConvertFilter(s, sep string, convertFn func(s string) string, filterFn func(s string) bool) []string {
	tokens := strings.Split(s, sep)
	res := make([]string, 0, len(tokens))
	for _, token := range tokens {
		word := token
		if convertFn != nil {
			word = convertFn(word)
		}

		if filterFn == nil || filterFn(token) {
			res = append(res, word)
		}
	}
	return res
}
