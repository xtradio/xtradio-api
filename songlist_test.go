package main

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetSongsFromDBSuccess(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()

	rows := sqlmock.NewRows([]string{"id", "filename", "artist", "title", "album", "lenght", "share", "url", "image"}).
		AddRow(1, "foobar.mp3", "foo", "bar", "foo-bar", 10, "http://test.com/foo-bar", "foo-bar", "foo.jpg").
		AddRow(2, "barbaz.mp3", "bar", "baz", "bar-baz", 10, "http://test.com/bar-baz", "bar-baz", "bar.jpg")

	mock.ExpectQuery("SELECT id, filename, artist, title, album, lenght, share, url, image FROM details ORDER BY id DESC").WillReturnRows(rows)
	mock.ExpectCommit()

	getData, _ := getSongsFromDB(db)

	if len(getData) != 2 {
		t.Errorf("Expected a list with 2 items, got %d", len(getData))
	}

	for k, v := range getData {
		if k == 0 {
			if v.Image != "https://img.xtcd.in/tracks/foo.jpg" {
				t.Errorf("Expected https://img.xtcd.in/xtrack/foo.jpg, got %s", v.Image)
			}
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetSongsFromDBError(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectQuery("SELECT id, filename, artist, title, album, lenght, share, url, image FROM details ORDER BY id DESC").WillReturnError(fmt.Errorf("some error"))

	_, err = getSongsFromDB(db)

	if err == nil {
		t.Errorf("Was expecting an error, got none.")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
