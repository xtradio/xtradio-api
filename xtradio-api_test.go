package main

import "testing"

func TestHandleSongDetailsProperSongNames(t *testing.T) {
	artist := "test_artist"
	title := "test_title"
	filename := "test_filename"

	expectShow := "backup"

	gotArtist, gotTitle, gotFilename, gotShow := handleSongDetails(artist, title, filename)

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
}

func TestHandleSongDetailsNoArtistBackup(t *testing.T) {
	artist := ""
	title := "test_title"
	filename := "test_filename"

	expectShow := "backup"

	gotArtist, gotTitle, _, gotShow := handleSongDetails(artist, title, filename)

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

	expectArtist := "test_artist"
	expectTitle := "test_title"
	expectShow := "backup"

	gotArtist, gotTitle, _, gotShow := handleSongDetails(artist, title, filename)

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

	expectArtist := "test_artist"
	expectTitle := "test_title"
	expectShow := "live"

	gotArtist, gotTitle, _, gotShow := handleSongDetails(artist, title, filename)

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
