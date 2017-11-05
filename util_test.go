package main

import (
	"os"
	"testing"
)

func TestIsExit(t *testing.T) {
	pwd, _ := os.Getwd()

	file := pwd + "/util_test.go"
	if !isExist(file) {
		t.Errorf("Expect isExist returns true but false. file: %s", file)
	}

	file = pwd + "/unexist.go"
	if isExist(file) {
		t.Errorf("Expect isExist returns false but true. file: %s", file)
	}
}

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

func TestIsBlank(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{input: " ", want: true},
		{input: "  ", want: true},
		{input: "\t", want: true},
		{input: "a", want: false},
		{input: "　", want: false},
	}

	for _, test := range tests {
		result := isBlank(test.input)

		if test.want != result {
			t.Errorf("expected %v, but %v. input: %v", test.want, result, test.input)
		}
	}
}
