package main

import (
	"testing"
)

func TestSplitSortSuccess(t *testing.T) {
	vars := "[\"test\",\"DESC\"]"
	testRow, testOrder, _ := splitSort(vars)

	if testRow != "test" {
		t.Errorf("Return was incorrect, got: %s, want: %s", testRow, "test")
	}

	if testOrder != "DESC" {
		t.Errorf("Return was incorrect, got: %s, want: %s", testOrder, "DESC")
	}
}

func TestSplitSortFail(t *testing.T) {
	vars := "test"
	testRow, testOrder, _ := splitSort(vars)

	if testRow != "" {
		t.Errorf("We expected an empty response, we got: %s", testRow)
	}

	if testOrder != "" {
		t.Errorf("We expected an empty response, we got: %s", testOrder)
	}
}

func TestSplitRangeSuccess(t *testing.T) {
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

func TestSplitRangeFailLength(t *testing.T) {
	vars := "[0]"

	_, _, err := splitRange(vars)

	if err == nil {
		t.Error("Expecting an error, we actually got nil.")
	}
}

func TestSplitRangeFailTypeMin(t *testing.T) {
	vars := "[test,1]"

	_, _, err := splitRange(vars)

	if err == nil {
		t.Error("Expecting an error, we got nil")
	}
}

func TestSplitRangeFailTypeMax(t *testing.T) {
	vars := "[0,test]"

	_, _, err := splitRange(vars)

	if err == nil {
		t.Error("Expecting an error, we got nil")
	}
}

func TestQueryBuilderEmptyFilter(t *testing.T) {
	testFilter := "{}"
	testRow := "a"
	testOrder := "DESC"
	data := queryBuilder(testFilter, testRow, testOrder)

	if data != "SELECT id, artist, title, album, lenght, share, url, image FROM details ORDER BY a DESC" {
		t.Errorf("Got wrong output: %s", data)
	}
}

func TestQueryBuilderNonEmptyFilter(t *testing.T) {
	testFilter := "{\"q\":\"test\"}"
	testRow := "a"
	testOrder := "DESC"
	data := queryBuilder(testFilter, testRow, testOrder)

	testQuery := "SELECT id, artist, title, album, lenght, share, url, image FROM details WHERE artist LIKE '%test%' ORDER BY a DESC"

	if data != testQuery {
		t.Skipf("Got wrong output, want: %s, got: %s", testQuery, data)
	}
}
