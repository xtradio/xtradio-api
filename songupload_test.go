package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestDoesExist(t *testing.T) {
	testFile := "test.txt"
	getResult := doesExist(testFile)

	if getResult != false {
		t.Errorf("Expected false, got %t", getResult)
	}

	testEmptyFile, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Error creating file for testing: %s", err)
	}
	testEmptyFile.Close()

	getResult = doesExist(testFile)

	if getResult != true {
		t.Errorf("Expected true, got %t", getResult)
	}

	err = os.Remove(testFile)
	if err != nil {
		t.Fatalf("Error removing %s for cleanup", testFile)
	}
}

func TestDuration(t *testing.T) {
	testMp3File := "test.mp3"

	_, err := duration(testMp3File)

	if err == nil {
		t.Errorf("Expected the file to not exist and an error, did not get an error.")
	}
}

func TestSaveDataSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testFilename := "test/test.mp3"
	testArtist := "test"
	testTitle := "foo"
	testDuration := float64(10)
	testShare := "https://test.com/url"
	testImage := "test.png"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO details").WithArgs(testFilename, testArtist, testTitle, "", testDuration, testShare, "", testImage, "daily", 0, 1).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	getResult, err := saveData(db, testFilename, testArtist, testTitle, testDuration, testShare, testImage)

	if getResult != 1 {
		t.Errorf("Expected id returned 1, got %d", getResult)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestSaveDataFail(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	testFilename := "test/test.mp3"
	testArtist := "test"
	testTitle := "foo"
	testDuration := float64(10)
	testShare := "https://test.com/url"
	testImage := "test.png"

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO details").WithArgs(testFilename, testArtist, testTitle, "", testDuration, testShare, "", testImage, "daily", 0, 1).WillReturnError(fmt.Errorf("some error"))
	mock.ExpectCommit()

	_, err = saveData(db, testFilename, testArtist, testTitle, testDuration, testShare, testImage)

	if err == nil {
		t.Errorf("Was expecting an error, got none.")
	}
}
