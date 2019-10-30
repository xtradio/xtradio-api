package main

import (
	"encoding/json"
	"fmt"
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
	sort0 := strings.Trim(sort[0], "[")
	sort1 := strings.Trim(sort[1], "]")
	row = strings.Replace(sort0, "\"", "", -1)
	order = strings.Replace(sort1, `"`, ``, -1)

	return row, order, nil
}

func splitRange(vars string) (min int64, max int64, err error) {
	rangeminmax := strings.Split(vars, ",")
	min, err = strconv.ParseInt(strings.Trim(rangeminmax[0], "["), 10, 64)
	if err != nil {
		fmt.Println("Error changing min string to int64 ", err)
		return
	}
	max, err = strconv.ParseInt(strings.Trim(rangeminmax[1], "]"), 10, 64)
	if err != nil {
		fmt.Println("Error changing max string to int64 ", err)
		return
	}

	return min, max, nil
}

func songList(w http.ResponseWriter, r *http.Request) {
	vars := r.URL.Query()
	var v struct {
		Data  []SongDetails `json:"data"`
		Total int64         `json:"total"`
	}

	var (
		count int64
		query string
	)

	row, order, err := splitSort(vars["sort"][0])

	min, max, err := splitRange(vars["range"][0])

	if vars["filter"][0] != "{}" {
		searchQuery1 := strings.Split(vars["filter"][0], ":")
		searchQuery2 := strings.Trim(searchQuery1[1], "}")
		searchQuery := "'%" + strings.Replace(searchQuery2, "\"", "", -1) + "%'"
		query = fmt.Sprintf("SELECT id, artist, title, album, lenght, share, url, image FROM details WHERE artist LIKE %s ORDER BY %s %s ", searchQuery, row, order)
	} else {
		query = fmt.Sprintf("SELECT id, artist, title, album, lenght, share, url, image FROM details ORDER BY %s %s", row, order)
	}

	fmt.Println(query)
	// Fetch details for the track
	rows := querysql(query)

	count = 0

	for rows.Next() {
		var s SongDetails

		err := rows.Scan(&s.ID, &s.Artist, &s.Title, &s.Album, &s.Length, &s.Share, &s.URL, &s.Image)
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
