package main

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestHandleSongDetailsProperSongNames(t *testing.T) {
	artist := "test_artist"
	title := "test_title"
	filename := "test_filename"
	length := "33"
	share := "http://foo.bar"
	image := "default.png"

	expectShow := "backup"

	gotArtist, gotTitle, gotFilename, gotShow, gotImage, gotShare, gotLength := handleSongDetails(artist, title, filename, image, share, length)

	if gotArtist != artist {
		t.Errorf("Expected %s, got %s", artist, gotArtist)
	}

	if gotTitle != title {
		t.Errorf("Expected %s, got %s", title, gotTitle)
	}

	if gotFilename != filename {
		t.Errorf("Expected %s, got %s", filename, gotFilename)
	}

	if gotShow != expectShow {
		t.Errorf("Expected %s, got %s", expectShow, gotShow)
	}

	if gotImage != image {
		t.Errorf("Expected %s, got %s", image, gotImage)
	}

	if gotShare != share {
		t.Errorf("Expected %s, got %s", gotShare, share)
	}

	if reflect.TypeOf(gotLength).Kind() != reflect.Float64 {
		t.Errorf("Expected float64, got %s", reflect.TypeOf(gotLength).Kind())
	}
}

func TestHandleSongDetailsNoArtistBackup(t *testing.T) {
	artist := ""
	title := "test_title"
	filename := "test_filename"
	length := "33"
	share := "http://foo.bar"
	image := "default.png"

	expectShow := "backup"

	gotArtist, gotTitle, _, gotShow, _, _, _ := handleSongDetails(artist, title, filename, image, share, length)

	if gotArtist != artist {
		t.Errorf("Expected %s, got %s", artist, gotArtist)
	}

	if gotTitle != title {
		t.Errorf("Expected %s, got %s", title, gotTitle)
	}

	if gotShow != expectShow {
		t.Errorf("Expected %s, got %s", expectShow, gotShow)
	}
}

func TestHandleSongDetailsMissingArtistBackup(t *testing.T) {
	artist := ""
	title := "test_artist - test_title"
	filename := "test_filename"
	length := "33"
	share := "http://foo.bar"
	image := "default.png"

	expectArtist := "test_artist"
	expectTitle := "test_title"
	expectShow := "backup"

	gotArtist, gotTitle, _, gotShow, _, _, _ := handleSongDetails(artist, title, filename, image, share, length)

	if gotArtist != expectArtist {
		t.Errorf("Expected %s, got %s", expectArtist, gotArtist)
	}

	if gotTitle != expectTitle {
		t.Errorf("Expected %s, got %s", expectTitle, gotTitle)
	}

	if gotShow != expectShow {
		t.Errorf("Expected %s, got %s", expectShow, gotShow)
	}
}

func TestHandleSongDetailsMissingArtistLiveshow(t *testing.T) {
	artist := ""
	title := "test_artist - test_title / Live DJ"
	filename := ""
	length := "33"
	share := "http://foo.bar"
	image := "default.png"

	expectArtist := "test_artist"
	expectTitle := "test_title"
	expectShow := "live"

	gotArtist, gotTitle, _, gotShow, _, _, _ := handleSongDetails(artist, title, filename, image, share, length)

	if gotArtist != expectArtist {
		t.Errorf("Expected %s, got %s", expectArtist, gotArtist)
	}

	if gotTitle != expectTitle {
		t.Errorf("Expected %s, got %s", expectTitle, gotTitle)
	}

	if gotShow != expectShow {
		t.Errorf("Expected %s, got %s", expectShow, gotShow)
	}
}

func TestLogSongToDBDisabledAndNoDB(t *testing.T) {
	os.Setenv("LOG_TO_DB", "false")

	var testSong Song
	testFilename := "/path/to/test/file.mp3"
	db, _, _ := sqlmock.New()
	defer db.Close()

	testID, _ := logSongToDB(db, testSong, testFilename)

	if testID != 0 {
		t.Errorf("Exepcted 0 id, got %d", testID)
	}

	os.Setenv("LOG_TO_DB", "")

}

func TestLogSongToDBSuccess(t *testing.T) {
	os.Setenv("LOG_TO_DB", "true")
	var testSong Song
	testFilename := "/path/to/test/file.mp3"
	testSong.Artist = "testArtist"
	testSong.Title = "testTitle"
	testSong.Show = "testShow"
	testSong.Image = "test.jpg"
	testSong.Share = "https://test.com"
	testSong.Length = 300

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectPrepare("INSERT INTO playlist")
	mock.ExpectExec("INSERT INTO playlist").WithArgs(testSong.Artist, testSong.Title, testFilename, testSong.Title, time.Now().Local().Format("2006-01-02"), time.Now().Local().Format("15:04:05")).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	testID, err := logSongToDB(db, testSong, testFilename)

	if err != nil {
		t.Errorf("Expected nil error, got %s", err)
	}

	if testID != 1 {
		t.Errorf("Expecting 1, got %d", testID)
	}

	// we make sure that all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

	os.Setenv("LOG_TO_DB", "")

}
