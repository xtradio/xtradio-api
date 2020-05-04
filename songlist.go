package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// SongDetails to output details of the songs to json
type SongDetails struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Artist   string `json:"artist"`
	Show     string `json:"show"`
	Image    string `json:"image"`
	Filename string `json:"filename"`
	Album    string `json:"album"`
	Length   string `json:"lenght"`
	Share    string `json:"share"`
	URL      string `json:"url"`
}

func songList(w http.ResponseWriter, r *http.Request) {
	log.Printf("songList function called by %s", r.RemoteAddr)
	var v struct {
		Data []SongDetails `json:"data"`
	}

	db, err := dbConnection()
	if err != nil {
		log.Println("Error opening connection to DB: ", err)
		return
	}
	defer db.Close()

	v.Data, err = getSongsFromDB(db)
	if err != nil {
		log.Println("Error getting items from database: ", err)
		return
	}

	json.NewEncoder(w).Encode(v)
}

func getSongsFromDB(db *sql.DB) ([]SongDetails, error) {
	var l []SongDetails

	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	rows, err := db.Query("SELECT id, filename, artist, title, album, lenght, share, url, image FROM details ORDER BY id DESC")

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var s SongDetails

		rows.Scan(&s.ID, &s.Filename, &s.Artist, &s.Title, &s.Album, &s.Length, &s.Share, &s.URL, &s.Image)

		if s.Image == "" {
			s.Image = "default.png"
		}
		s.Image = fmt.Sprintf("https://img.xtcd.in/tracks/%s", s.Image)

		l = append(l, s)

	}

	defer rows.Close()

	return l, nil
}
