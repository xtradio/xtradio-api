package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Song details
type Song struct {
	Title     string `json:"song"`
	Artist    string `json:"artist"`
	Show      string `json:"show"`
	Image     string `json:"image"`
	Album     string `json:"album"`
	Length    int    `json:"length"`
	Remaining int    `json:"remaining"`
	Share     string `json:"share"`
	URL       string `json:"url"`
}

// Duration of the song
type Duration struct {
	Finished time.Time
}

type cache struct {
	sync.RWMutex
	song     Song
	duration Duration
}

type songsHandler struct {
	c *cache
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XTRadio API.")
	fmt.Println("Endpoint Hit: homePage")
}

func (h songsHandler) readPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var song Song
	var duration Duration

	fmt.Println(time.Now(), r.RequestURI, "Reading post message for song.", vars["file"])

	// Open and connect do DB
	db, err := sql.Open("mysql", "root:test@tcp(db:3306)/radio?charset=utf8")
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.")
		return
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		http.Error(w, err.Error(), 500)
		fmt.Println("Ping database failed.")
		return
	}

	// Fetch details for the track
	query := db.QueryRow("SELECT artist, title, album, lenght, share, url, image FROM details WHERE filename=?", vars["file"])

	// Save it to the "Songs" struct
	err = query.Scan(&song.Artist, &song.Title, &song.Album, &song.Length, &song.Share, &song.URL, &song.Image)
	if err != nil {
		// If the file is not found in the db split it
		artist := strings.Split(vars["file"], "/")
		// Return the last element
		song.Artist = artist[len(artist)-1]
		// Replace _
		song.Artist = strings.Replace(song.Artist, "_", " ", -1)
		// Replace ".mp3"
		song.Artist = strings.Replace(song.Artist, ".mp3", "", -1)
		song.Title = ""
		song.Album = ""
		song.Length = 0
		song.Share = ""
		song.URL = ""

		fmt.Println(time.Now(), "Scan not found.")
	}

	if song.Image == "" {
		song.Image = "default.png"
	}

	// Add url to image
	song.Image = "https://img.xtradio.org/tracks/" + song.Image

	// Calculate when the song will finish, will be needed for the "remaining" var
	duration.Finished = time.Now().Local().Add(time.Second * time.Duration(song.Length))

	defer db.Close()

	h.c.Lock()
	defer h.c.Unlock()
	h.c.song = song
	h.c.duration = duration
}

func (h songsHandler) returnSongs(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()

	// Calculate remaining seconds in real time
	remaining := time.Until(h.c.duration.Finished)
	h.c.song.Remaining = int(remaining.Seconds())
	fmt.Println("Time remaining: ", remaining.Seconds())

	if remaining.Seconds() < 0 {
		fmt.Println(time.Now(), r.RequestURI, "Song duration expired.")
		http.Error(w, "API Unavailable", 503)
		return
	}

	fmt.Println(time.Now(), r.RequestURI, "Served api request.")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Output json
	json.NewEncoder(w).Encode(h.c.song)
}

func publishAPI() {
	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage)
	sh := songsHandler{c: &cache{}}
	apiRouter.HandleFunc("/api", sh.returnSongs)
	apiRouter.HandleFunc("/post/song", sh.readPost).
		Name("putsong").
		Queries("file", "{file}")

	log.Fatal(http.ListenAndServe(":10000", apiRouter))
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println("Rest API v2.0 - Mux Routers")
	publishAPI()
}
