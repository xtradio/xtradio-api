package main

import (
	"testing"
)

func TestSplitSort(t *testing.T) {
	vars := "[\"test\",\"DESC\"]"
	testRow, testOrder, err := splitSort(vars)

	if testRow != "test" {
		t.Errorf("Return was incorrect, got: %s, want: %s", testRow, "test")
	}

	if testOrder != "DESC" {
		t.Errorf("Return was incorrect, got: %s, want: %s", testOrder, "DESC")
	}

	if err != nil {
		t.Errorf("Error was returned")
	}
}

func TestSplitRange(t *testing.T) {
	vars := "[0,1]"

	testMin, testMax, err := splitRange(vars)

	if testMin != 0 {
		t.Errorf("Return was incorrect, got: %d, want: 0", testMin)
	}

	if testMax != 1 {
		t.Errorf("Return was incorrect, got: %d, want: 1", testMax)
	}

	if err != nil {
		t.Errorf("Error was returned")
	}
}
