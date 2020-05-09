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
	"net/url"
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

	f.Artist, f.Title, f.URL, f.Image = r.FormValue("artist"), r.FormValue("title"), r.FormValue("url"), r.FormValue("image")

	db, err := dbConnection()
	defer db.Close()

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
	image, err := sendToCDN(f.Image)
	if err != nil {
		log.Println("Saving image failed: ", err)
		h.Response = "fail"
		h.Reason = "Failed to save image."
		json.NewEncoder(w).Encode(h)
		return
	}

	savedID, err := saveData(db, filename, f.Artist, f.Title, fileDuration, f.URL, image)
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

	directory, err := getEnv("MUSIC_UPLOAD_DIR")
	if err != nil {
		return "", err
	}

	filename := fmt.Sprintf("%s/%s", directory, header.Filename)

	ok := doesExist(filename)
	if ok != false {
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
		return false
	}
	return true
}

func saveData(db *sql.DB, filename string, artist string, title string, duration float64, share string, image string) (int64, error) {
	// Add raw data in to database
	// insert

	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()

	stmt, err := tx.Exec("INSERT INTO details (filename, artist, title, album, lenght, share, url, image, playlist, vote, review) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", filename, artist, title, "", duration, share, "", image, "daily", 0, 1)
	if err != nil {
		return 0, err
	}

	id, err := stmt.LastInsertId()
	if err != nil {
		return 0, err
	}

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

func sendToCDN(imgURL string) (string, error) {
	type cdnResponse struct {
		Response string `json:"response"`
		Filename string `json:"filename"`
	}

	var URL *url.URL
	var c cdnResponse

	cdnURL, err := getEnv("CDN_URL")
	if err != nil {
		return "", err
	}

	URL, _ = url.Parse(fmt.Sprintf("http://%s/v1/upload", cdnURL))

	data := url.Values{}
	data.Set("imgURL", imgURL)

	URL.RawQuery = data.Encode()

	res, err := http.Get(URL.String())
	if err != nil {
		return "", err
	}

	if res.StatusCode != 200 {
		errMsg := fmt.Errorf("response code from URL was %d", res.StatusCode)
		return "", errMsg
	}

	f, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	res.Body.Close()

	json.Unmarshal(f, &c)

	if c.Response != "ok" {
		errMsg := fmt.Errorf("failed to upload image to CDN")
		return "", errMsg
	}

	return c.Filename, nil
}
