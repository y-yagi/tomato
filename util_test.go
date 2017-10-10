package main

import "testing"

func TestContains(t *testing.T) {
	value := "today"
	list := []string{"today", "week", "month", "all"}

	if !contains(list, value) {
		t.Errorf("Expect contains returns true but false. list: %v, value: '%s'", list, value)
	}

	value = "tomorrow"
	if contains(list, value) {
		t.Errorf("Expect contains returns false but true. list: %v, value: '%s'", list, value)
	}
}
