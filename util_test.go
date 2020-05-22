package main

import (
	"os"
	"testing"
)

func TestGetEnvTrue(t *testing.T) {
	testEnvVarValue := "testing"
	testEnvVarKey := "TEST_ENV"
	os.Setenv(testEnvVarKey, testEnvVarValue)

	returnEnv := getEnv(testEnvVarKey)

	if returnEnv != testEnvVarValue {
		t.Errorf("Return was incorect, got: %s, want: %s", returnEnv, testEnvVarValue)
	}

}

func TestGetEnvFalse(t *testing.T) {
	testReturn := getEnv("testing")

	if testReturn != "" {
		t.Error("Value returned, we expected an error")
	}

}

func TestSongHistory(t *testing.T) {
	var songList []Song

	newSong := Song{
		Artist:    "new",
		Title:     "test",
		Show:      "backup",
		Image:     "https://img.xtradio.org/tracks/default.png",
		Album:     "",
		Length:    123,
		Remaining: 5,
		Share:     "http://soundcloud.com/xtradio",
	}

	existingSong := Song{
		Artist:    "first",
		Title:     "exist",
		Show:      "backup",
		Image:     "https://img.xtradio.org/tracks/default.png",
		Album:     "",
		Length:    123,
		Remaining: 5,
		Share:     "http://soundcloud.com/xtradio",
	}

	existingSong1 := Song{
		Artist:    "second",
		Title:     "exist",
		Show:      "backup",
		Image:     "https://img.xtradio.org/tracks/default.png",
		Album:     "",
		Length:    123,
		Remaining: 5,
		Share:     "http://soundcloud.com/xtradio",
	}

	existingSong2 := Song{
		Artist:    "third",
		Title:     "Song",
		Show:      "backup",
		Image:     "https://img.xtradio.org/tracks/default.png",
		Album:     "",
		Length:    123,
		Remaining: 5,
		Share:     "http://soundcloud.com/xtradio",
	}

	getData := songHistory(songList, newSong)

	if len(getData) != 1 {
		t.Errorf("Failed for number of items, expected 1 item, got %d", len(getData))
	}

	for _, v := range getData {
		if v.Artist != newSong.Artist {
			t.Errorf("Failed for one item, expected %s, got %s: ", newSong.Artist, v.Artist)
		}
	}

	// reset list
	songList = nil

	// Test for two songs in the list
	songList = append(songList, existingSong)
	getData = songHistory(songList, newSong)

	if len(getData) != 2 {
		t.Errorf("Failed for number of items, expected 2 items, got %d", len(getData))
	}

	for k, v := range getData {
		if k == 0 && v.Artist != newSong.Artist {
			t.Errorf("Expecting two items: %s, got %s", newSong.Artist, v.Artist)
		} else if k == 1 && v.Artist != existingSong.Artist {
			t.Errorf("Expecting two items: %s, got %s", existingSong.Artist, v.Artist)
		}
	}

	// reset list
	songList = nil

	// Test for three songs
	songList = append(songList, existingSong)
	songList = append(songList, existingSong1)
	getData = songHistory(songList, newSong)

	if len(getData) != 3 {
		t.Errorf("Failed on number of items, expecting 3 items, got %d", len(getData))
	}

	for k, v := range getData {
		if k == 0 && v.Artist != newSong.Artist {
			t.Errorf("3 song test: Expecting 3 items: %s, got %s", newSong.Artist, v.Artist)
		} else if k == 1 && v.Artist != existingSong.Artist {
			t.Errorf("3 song test: Expecting 3 items: %s, got %s", existingSong.Artist, v.Artist)
		} else if k == 2 && v.Artist != existingSong1.Artist {
			t.Errorf("3 song test: Expecting 3 items: %s, got %s", existingSong1.Artist, v.Artist)
		}
	}

	// reset list
	songList = nil

	// Test for four songs
	songList = append(songList, existingSong2)
	songList = append(songList, existingSong1)
	songList = append(songList, existingSong)
	getData = songHistory(songList, newSong)

	if len(getData) != 3 {
		t.Errorf("We expect only 3 items in the slice, we got %d", len(getData))
	}

	for k, v := range getData {
		if k == 0 && v.Artist != newSong.Artist {
			t.Errorf("4 song test: Expecting 3 items: %s, got %s", newSong.Artist, v.Artist)
		} else if k == 1 && v.Artist != existingSong2.Artist {
			t.Errorf("4 song test: Expecting 3 items: %s, got %s", existingSong2.Artist, v.Artist)
		} else if k == 2 && v.Artist != existingSong1.Artist {
			t.Errorf("4 song test: Expecting 3 items: %s, got %s", existingSong1.Artist, v.Artist)
		}
	}

}
