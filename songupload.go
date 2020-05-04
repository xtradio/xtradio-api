package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/tcolgate/mp3"
)

func songUpload(w http.ResponseWriter, r *http.Request) {

	type formDetails struct {
		Artist   string
		Title    string
		Duration float64
		URL      string
		Image    string
	}

	type httpResponse struct {
		Response string `json:"response"`
		Reason   string `json:"reason"`
	}

	var f formDetails
	var h httpResponse

	f.Artist = r.FormValue("artist")
	f.Title = r.FormValue("title")
	// f.Duration = r.FormValue("duration/")
	f.URL = r.FormValue("url")
	f.Image = r.FormValue("image")

	log.Println(f)

	filename, err := saveFile(r.FormFile("file"))
	if err != nil {
		log.Println("Saving mp3 file failed:", err)
		h.Response = "fail"
		h.Reason = fmt.Sprintf("Failed to save mp3 file: %s", err)
		json.NewEncoder(w).Encode(h)
		return
	}

	fileDuration, err := duration(filename)
	if err != nil {
		log.Println(err)
		h.Response = "fail"
		h.Reason = "Failed to get duration of uploaded file."
		json.NewEncoder(w).Encode(h)
		return
	}

	// TODO: Send image to CDN to save it.
	image := ""
	// image, err := downloadFile("files/test.jpg", f.Image)
	// if err != nil {
	// 	log.Println("Saving image failed: ", err)
	// 	h.Response = "fail"
	// 	h.Reason = "Failed to save image."
	// 	json.NewEncoder(w).Encode(h)
	// 	return
	// }

	savedID, err := saveData(filename, f.Artist, f.Title, fileDuration, f.URL, image)
	if err != nil {
		log.Println("Saving data to DB failed:", err)
		h.Response = "fail"
		h.Reason = "Failed to save data in to database."
		json.NewEncoder(w).Encode(h)
		return
	}

	h.Response = "success"
	h.Reason = fmt.Sprintf("Saved data with id %d", savedID)

	// do something else
	// etc write header
	json.NewEncoder(w).Encode(h)
	return
}

func saveFile(file multipart.File, header *multipart.FileHeader, _ error) (string, error) {
	var Buf bytes.Buffer

	defer file.Close()

	directory, err := getEnv("XTRADIO_MUSIC_UPLOAD_DIR")
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/%s", directory, header.Filename)

	ok := doesExist(filename)
	if ok != true {
		return "", fmt.Errorf("File %s already exists", filename)
	}
	// Copy the file data to my buffer
	io.Copy(&Buf, file)

	err = ioutil.WriteFile(filename, Buf.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	Buf.Reset()

	return filename, nil
}

func doesExist(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return true
	}
	return false
}

func saveData(filename string, artist string, title string, duration float64, share string, image string) (int64, error) {
	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Println("Ping database failed.", err)
		return 0, err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Println("Ping database failed.", err)
		return 0, err
	}

	// Add raw data in to database
	// insert
	stmt, err := db.Prepare("INSERT INTO details (filename, artist, title, album, lenght, share, url, image, playlist, vote, review) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Prepare of SQL statement failed", err)
		return 0, err
	}

	res, err := stmt.Exec(filename, artist, title, "", duration, share, "", image, "daily", 0, 1)
	if err != nil {
		log.Println("Adding data in to details table failed", err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error fetching last inserted ID", err)
		return 0, err
	}

	log.Println("Inserted last played song with id: ", id)
	return id, nil
}

func duration(file string) (float64, error) {
	t := 0.0

	r, err := os.Open(file)
	if err != nil {
		return t, err
	}

	d := mp3.NewDecoder(r)
	var f mp3.Frame
	skipped := 0

	for {

		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
			return t, err
		}

		t = t + f.Duration().Seconds()
	}

	return t, nil

}

// func downloadFile(filepath string, url string) (string, error) {

// 	// Get the data
// 	resp, err := http.Get(url)
// 	if err != nil {
// 		return filepath, err
// 	}
// 	defer resp.Body.Close()

// 	// Create the file
// 	out, err := os.Create(filepath)
// 	if err != nil {
// 		return filepath, err
// 	}
// 	defer out.Close()

// 	// Write the body to file
// 	_, err = io.Copy(out, resp.Body)
// 	return filepath, err
// }
