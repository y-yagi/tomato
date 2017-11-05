package main

import (
	"os"
	"regexp"
)

var blank_re = regexp.MustCompile("\\A[[:space:]]*\\z")

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func contains(s []string, e string) bool {
	for _, v := range s {
		if e == v {
			return true
		}
	}
	return false
}

func isBlank(str string) bool {
	return len(str) == 0 || blank_re.MatchString(str)
}
