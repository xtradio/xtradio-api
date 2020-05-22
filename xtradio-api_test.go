package main

import (
	"reflect"
	"testing"
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
