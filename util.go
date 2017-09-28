package main

import "os"

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
