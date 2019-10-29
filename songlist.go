package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func songList(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)
	vars := r.URL.Query()
	var v struct {
		Data  []SongDetails `json:"data"`
		Total int64         `json:"total"`
	}

	var (
		count int64
		query string
	)

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.", err)
		return
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.", err)
		return
	}

	sort := strings.Split(vars["sort"][0], ",")
	sort0 := strings.Trim(sort[0], "[")
	sort1 := strings.Trim(sort[1], "]")
	row := strings.Replace(sort0, "\"", "", -1)
	order := strings.Replace(sort1, `"`, ``, -1)

	rangeminmax := strings.Split(vars["range"][0], ",")
	min, err := strconv.ParseInt(strings.Trim(rangeminmax[0], "["), 10, 64)
	if err != nil {
		fmt.Println("Error changing min string to int64 ", err)
		return
	}
	max, err := strconv.ParseInt(strings.Trim(rangeminmax[1], "]"), 10, 64)
	if err != nil {
		fmt.Println("Error changing max string to int64 ", err)
		return
	}

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
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("Fetching rows failed.", err)
		return
	}

	defer rows.Close()

	count = 0

	for rows.Next() {
		var s SongDetails

		err := rows.Scan(&s.Id, &s.Artist, &s.Title, &s.Album, &s.Length, &s.Share, &s.URL, &s.Image)
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

}
