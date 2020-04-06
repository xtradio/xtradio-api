package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// SongDetails to output details of the songs to json
type SongDetails struct {
	ID     int64  `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Show   string `json:"show"`
	Image  string `json:"image"`
	Album  string `json:"album"`
	Length string `json:"lenght"`
	Share  string `json:"share"`
	URL    string `json:"url"`
}

func splitSort(vars string) (row string, order string, err error) {

	sort := strings.Split(vars, ",")
	if len(sort) != 2 {
		err = errors.New("not the expected amount of parameters for sort string")
		return
	}
	sort0 := strings.Trim(sort[0], "[")
	sort1 := strings.Trim(sort[1], "]")
	row = strings.Replace(sort0, "\"", "", -1)
	order = strings.Replace(sort1, `"`, ``, -1)

	return row, order, nil
}

func splitRange(vars string) (min int64, max int64, err error) {
	rangeminmax := strings.Split(vars, ",")

	if len(rangeminmax) != 2 {
		err = errors.New("not the expected amount of parameters for range string")
		return
	}

	min, err = strconv.ParseInt(strings.Trim(rangeminmax[0], "["), 10, 64)
	if err != nil {
		log.Printf("Error changing min string to int64: %s", err)
		return
	}
	max, err = strconv.ParseInt(strings.Trim(rangeminmax[1], "]"), 10, 64)
	if err != nil {
		log.Printf("Error changing max string to int64: %s ", err)
		return
	}

	return min, max, nil
}

func queryBuilder(filter string, row string, order string) (query string) {
	if filter == "{}" {
		query = fmt.Sprintf("SELECT id, artist, title, album, lenght, share, url, image FROM details ORDER BY %s %s", row, order)
		return query
	}
	searchQuery1 := strings.Split(filter, ":")
	searchQuery2 := strings.Trim(searchQuery1[1], "}")
	searchQuery := "'%" + strings.Replace(searchQuery2, "\"", "", -1) + "%'"

	query = fmt.Sprintf("SELECT id, artist, title, album, lenght, share, url, image FROM details WHERE artist LIKE %s ORDER BY %s %s ", searchQuery, row, order)

	return query
}

func songList(w http.ResponseWriter, r *http.Request) {
	log.Printf("songList function called by %s", r.RemoteAddr)
	vars := r.URL.Query()
	var v struct {
		Data  []SongDetails `json:"data"`
		Total int64         `json:"total"`
	}

	var (
		count int64
	)

	row, order, err := splitSort(vars["sort"][0])

	min, max, err := splitRange(vars["range"][0])

	query := queryBuilder(vars["filter"][0], row, order)

	rows, err := querysql(query)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	count = 0

	for rows.Next() {
		var s SongDetails

		err := rows.Scan(&s.ID, &s.Artist, &s.Title, &s.Album, &s.Length, &s.Share, &s.URL, &s.Image)
		s.Image = "https://img.xtradio.org/tracks/" + s.Image
		if err != nil {
			http.Error(w, err.Error(), 500)
			fmt.Println("Fetching item failed.", err)
			return
		}
		if (count >= min) && (count <= max) {
			v.Data = append(v.Data, s)
		}
		count = count + 1
	}
	v.Total = count
	p, err := json.Marshal(v)
	if err != nil {
		fmt.Println("Error on defining json", err)
		return
	}
	fmt.Fprintf(w, string(p))

	defer rows.Close()

}
