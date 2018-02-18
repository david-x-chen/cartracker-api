package main

import "strings"

func stringInSlice(a string, list []string) bool {
	var result bool
	for _, b := range list {
		result = strings.EqualFold(b, a)
	}
	return result
}
