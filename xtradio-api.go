package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/alexandrevicenzi/go-sse"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Song details
type Song struct {
	Title     string  `json:"song"`
	Artist    string  `json:"artist"`
	Show      string  `json:"show"`
	Image     string  `json:"image"`
	Album     string  `json:"album"`
	Length    float64 `json:"length"`
	Remaining float64 `json:"remaining"`
	Share     string  `json:"share"`
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

func handleSongDetails(artist string, title string, filename string, image string, share string, length string) (string, string, string, string, string, string, float64) {
	var show string

	if image == "" {
		image = "default.png"
	}

	lengthFloat64, err := strconv.ParseFloat(length, 10)
	if err != nil {
		log.Println("Error converting string to float64", err)
	}

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
			return artist, title, filename, show, image, share, lengthFloat64
		}
	}

	show = "backup"

	return artist, title, filename, show, image, share, lengthFloat64
}

func logSongToDB(db *sql.DB, song Song, filename string) (int64, error) {
	if getEnv("LOG_TO_DB") != "true" {
		return 0, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return 0, nil
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit()
		default:
			tx.Rollback()
		}
	}()
	// Add raw data in to database
	// insert
	stmt, err := tx.Prepare("INSERT INTO playlist (artist, title, filename, song, datum, time) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Println("Prepare of SQL statement failed", err)
		return 0, err
	}

	res, err := stmt.Exec(song.Artist, song.Title, filename, song.Title, time.Now().Local().Format("2006-01-02"), time.Now().Local().Format("15:04:05"))
	if err != nil {
		log.Println("Adding data in to playlist failed", err)
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

func (h songsHandler) readPost(w http.ResponseWriter, r *http.Request, s *sse.Server) {
	vars := mux.Vars(r)
	var song Song
	var duration Duration

	h.c.previousData = songHistory(h.c.previousData, h.c.song)

	var filename string

	song.Artist, song.Title, filename, song.Show, song.Image, song.Share, song.Length = handleSongDetails(vars["artist"], vars["title"], vars["file"], vars["image"], vars["share"], vars["length"])

	log.Println(r.RequestURI, "Reading post message for song.", vars["file"])

	db, err := dbConnection()
	if err != nil {
		log.Println("Connection to DB failed.")
	} else {
		_, err := logSongToDB(db, song, filename)
		if err != nil {
			log.Println("Error inserting np song to DB.", err)
		}
	}
	// Add url to image
	song.Image = "https://img.xtcd.in/tracks/" + song.Image

	// Calculate when the song will finish, will be needed for the "remaining" var
	duration.Finished = time.Now().Local().Add(time.Second * time.Duration(song.Length))

	// defer db.Close()

	// sendTweet("♪ #np " + song.Artist + " - " + song.Title + " " + song.Share)
	tuneinAPI(song.Artist, song.Title)

	songsPlayed.With(prometheus.Labels{"artist": song.Artist, "title": song.Title, "show": song.Show}).Inc()

	h.c.Lock()
	defer h.c.Unlock()
	h.c.song = song
	h.c.duration = duration

	h.c.upcomingData, err = upcomingSongs()
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	sseOutput(s, h)

	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func (h songsHandler) returnSongs(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()

	// Calculate remaining seconds in real time
	remaining := time.Until(h.c.duration.Finished)
	h.c.song.Remaining = float64(remaining.Seconds())
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
	h.c.song.Remaining = float64(remaining.Seconds())

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
	w.Header().Set("Content-Type", "application/json")
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
	}).Queries("file", "{file}", "artist", "{artist}", "title", "{title}", "image", "{image}", "share", "{share}", "length", "{length}")
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
