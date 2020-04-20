package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	song         Song
	previousData []Song
	upcomingData []UpcomingSongs
	duration     Duration
}

// PreviousSongs is the struct to hold the list of songs previously played
type PreviousSongs []Song

// APIOutput is a struct to gather all types of data to be outputted
type APIOutput struct {
	CurrentSong   Song            `json:"current"`
	PreviousSongs PreviousSongs   `json:"previous"`
	UpcomingSongs []UpcomingSongs `json:"upcoming"`
}

type songsHandler struct {
	c *cache
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XTRadio API.")
	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func handleSongDetails(artist string, title string, filename string) (string, string, string, string) {
	var show string

	if artist == "" {
		splitTitle := strings.Split(title, " - ")

		if len(splitTitle) == 2 {

			artist = splitTitle[0]
			splitLive := strings.Split(splitTitle[1], " / ")
			title = splitLive[0]

			if len(splitLive) == 2 {
				if splitLive[1] == "Live DJ" {
					show = "live"
				}
			} else {
				show = "backup"
			}
			return artist, title, filename, show
		}
	}

	show = "backup"

	return artist, title, filename, show
}

func (h songsHandler) readPost(w http.ResponseWriter, r *http.Request, s *sse.Server) {
	vars := mux.Vars(r)
	var song Song
	var duration Duration

	h.c.previousData = songHistory(h.c.previousData, h.c.song)

	artist, title, filename, show := handleSongDetails(vars["artist"], vars["title"], vars["file"])

	log.Println(r.RequestURI, "Reading post message for song.", vars["file"])

	connection := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8", os.Getenv("MYSQL_USERNAME"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_HOST"), os.Getenv("MYSQL_DATABASE"))
	// Open and connect do DB
	db, err := sql.Open("mysql", connection)
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println("Ping database failed.", err)
		return
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		http.Error(w, err.Error(), 500)
		log.Println("Ping database failed.", err)
		return
	}

	// Add raw data in to database
	// insert
	stmt, err := db.Prepare("INSERT INTO playlist (artist, title, filename, song, datum, time) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Prepare of SQL statement failed", err)
		return
	}

	res, err := stmt.Exec(artist, title, filename, title, time.Now().Local().Format("2006-01-02"), time.Now().Local().Format("15:04:05"))
	if err != nil {
		log.Println("Adding data in to playlist failed", err)
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Println("Error fetching last inserted ID", err)
		return
	}

	log.Println("Inserted last played song with id: ", id)

	// Fetch details for the track
	query := db.QueryRow("SELECT artist, title, album, lenght, share, url, image FROM details WHERE filename=?", filename)

	// Save it to the "Songs" struct
	err = query.Scan(&song.Artist, &song.Title, &song.Album, &song.Length, &song.Share, &song.URL, &song.Image)
	if err != nil {
		// If the song is not in the db, use the metadata passed from liquidsoap
		song.Artist = artist
		song.Title = title
		song.Album = ""
		song.Length = 0
		song.Share = ""
		song.URL = ""
		log.Println(r.RemoteAddr, r.Method, r.URL, "Scan not found.")
	}

	if song.Image == "" {
		song.Image = "default.png"
	}

	song.Show = show

	// Add url to image
	song.Image = "https://img.xtradio.org/tracks/" + song.Image

	// Calculate when the song will finish, will be needed for the "remaining" var
	duration.Finished = time.Now().Local().Add(time.Second * time.Duration(song.Length))

	defer db.Close()

	// sendTweet("♪ #np " + song.Artist + " - " + song.Title + " " + song.Share)
	tuneinAPI(song.Artist, song.Title)

	h.c.Lock()
	defer h.c.Unlock()
	h.c.song = song
	h.c.duration = duration

	h.c.upcomingData = upcomingSongs()

	sseOutput(s, h)

	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func (h songsHandler) returnSongs(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()

	// Calculate remaining seconds in real time
	remaining := time.Until(h.c.duration.Finished)
	h.c.song.Remaining = int(remaining.Seconds() + 5)
	log.Println("Time remaining: ", remaining.Seconds())

	if remaining.Seconds() < 0 {
		log.Println(r.RemoteAddr, r.Method, r.URL, "Song duration expired - Faking time.")
		h.c.song.Remaining = 10
	}

	// Output json
	json.NewEncoder(w).Encode(h.c.song)
	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func nowplaying(h songsHandler) APIOutput {

	var data APIOutput

	// Calculate remaining seconds in real time
	remaining := time.Until(h.c.duration.Finished)
	h.c.song.Remaining = int(remaining.Seconds() + 1)

	if remaining.Seconds() < 0 {
		log.Println("Song duration expired - Faking time.")
		h.c.song.Remaining = 10
	}

	data.CurrentSong = h.c.song
	data.PreviousSongs = h.c.previousData
	data.UpcomingSongs = h.c.upcomingData

	return data
}

func (h songsHandler) apiOutput(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()

	data := nowplaying(h)

	// Output json
	json.NewEncoder(w).Encode(data)
	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func sseOutput(s *sse.Server, h songsHandler) {

	data := nowplaying(h)

	msg, err := json.Marshal(data)

	if err != nil {
		log.Println("There was an error creating the json blob, ", err)
	}

	s.SendMessage("", sse.SimpleMessage(string(msg)))
}

func publishAPI() {
	s := sse.NewServer(&sse.Options{
		// Increase default retry interval to 10m.
		RetryInterval: 10 * 6 * 10 * 1000,
		// Print debug info
		Logger: log.New(os.Stdout, "", log.Ldate|log.Ltime|log.Lshortfile),
	})
	defer s.Shutdown()

	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage)
	sh := songsHandler{c: &cache{}}
	apiRouter.HandleFunc("/api", sh.returnSongs)
	apiRouter.HandleFunc("/v1/np/", sh.apiOutput)
	apiRouter.HandleFunc("/post/song", func(w http.ResponseWriter, r *http.Request) {
		sh.readPost(w, r, s)
	}).Queries("file", "{file}", "artist", "{artist}", "title", "{title}")
	apiRouter.Handle("/v1/sse/np", s)
	// apiRouter.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	metricsServe := mux.NewRouter().StrictSlash(true)
	metricsServe.Handle("/metrics", promhttp.Handler())

	go func() {
		log.Fatal(http.ListenAndServe(":10001", metricsServe))
	}()

	log.Fatal(http.ListenAndServe(":10000", apiRouter))
}

func main() {
	log.Println("Rest API v2.0 - Mux Routers")
	publishAPI()
}
