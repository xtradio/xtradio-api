package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
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

// Status of the mountpoints
type Status struct {
	HighStatus string `json:"highstatus"`
	MidStatus  string `json:"midstatus"`
	LowStatus  string `json:"lowstatus"`
}

type cache struct {
	sync.RWMutex
	song     Song
	duration Duration
	status   Status
}

type songsHandler struct {
	c *cache
}

type statusHandler struct {
	c *cache
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XTRadio API.")
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL)
}

func (h statusHandler) readStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var status Status

	if vars["mount"] == "high" {
		status.HighStatus = vars["status"]
	} else if vars["mount"] == "mid" {
		status.MidStatus = vars["status"]
	} else if vars["mount"] == "low" {
		status.LowStatus = vars["status"]
	}

	fmt.Println(time.Now(), r.RequestURI, "Mount ", vars["mount"], " is ", vars["status"])

	h.c.status = status
}

func (h songsHandler) readPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var song Song
	var duration Duration

	fmt.Println(time.Now(), r.RequestURI, "Reading post message for song.", vars["file"])

	// Open and connect do DB
	db, err := sql.Open("mysql", "root:test@tcp(172.17.0.3:3306)/radio?charset=utf8")
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
		fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL, "Scan not found.")
	}

	if song.Image == "" {
		song.Image = "default.png"
	}

	// Add url to image
	song.Image = "https://img.xtradio.org/tracks/" + song.Image

	// Calculate when the song will finish, will be needed for the "remaining" var
	duration.Finished = time.Now().Local().Add(time.Second * time.Duration(song.Length))

	defer db.Close()

	sendTweet("♪ #np " + song.Artist + " - " + song.Title + " " + song.Share)
	tuneinAPI(song.Artist, song.Title)

	h.c.Lock()
	defer h.c.Unlock()
	h.c.song = song
	h.c.duration = duration
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL)
}

func tuneinAPI(artist string, title string) {

	partnerid := os.Getenv("TUNEIN_PARTNER_ID")
	partnerkey := os.Getenv("TUNEIN_PARTNER_KEY")
	stationid := os.Getenv("TUNEIN_STATION_ID")
	if partnerid == "" || partnerkey == "" || stationid == "" {
		fmt.Println(time.Now(), "No tunein creds, skipping.")
		return
	}

	var URL *url.URL
	URL, err := url.Parse("http://air.radiotime.com/Playing.ashx?")
	if err != nil {
		fmt.Println(time.Now(), "Tunein URL unavailable")
		return
	}

	parameters := url.Values{}
	parameters.Add("partnerId", partnerid)
	parameters.Add("partnerKey", partnerkey)
	parameters.Add("id", stationid)
	parameters.Add("artist", artist)
	parameters.Add("title", title)
	URL.RawQuery = parameters.Encode()

	fmt.Printf("Encoded URL is %q\n", URL.String())
	res, err := http.Get(URL.String())
	if err != nil {
		fmt.Println(err)
	}
	robots, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s", robots)
}

func sendTweet(message string) {
	consumerKey := os.Getenv("TWITTER_CONSUMER_KEY")
	consumerSecret := os.Getenv("TWITTER_CONSUMER_SECRET")
	accessToken := os.Getenv("TWITTER_ACCESS_TOKEN")
	accessSecret := os.Getenv("TWITTER_ACCESS_SECRET")
	if consumerKey == "" || consumerSecret == "" || accessToken == "" || accessSecret == "" {
		fmt.Println(time.Now(), "Missing required environment variable")
		return
	}
	config := oauth1.NewConfig(consumerKey, consumerSecret)
	token := oauth1.NewToken(accessToken, accessSecret)

	// httpClient will automatically authorize http.Request's
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)
	fmt.Println(time.Now(), "Tweet sent: "+message)
	tweet, resp, err := client.Statuses.Update(message, nil)
	if err != nil {
		fmt.Println(time.Now(), "Tweet not sent", tweet, resp, err)
	}
}

func (h statusHandler) returnStatus(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL, "Served status api request.")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	json.NewEncoder(w).Encode(h.c.status)
}

func (h songsHandler) returnSongs(w http.ResponseWriter, r *http.Request) {
	h.c.RLock()
	defer h.c.RUnlock()

	// Calculate remaining seconds in real time
	remaining := time.Until(h.c.duration.Finished)
	h.c.song.Remaining = int(remaining.Seconds())
	fmt.Println("Time remaining: ", remaining.Seconds())

	if remaining.Seconds() < 0 {
		fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL, "Song duration expired.")
		http.Error(w, "API Unavailable", 503)
		return
	}
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL, "Served api request.")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	// Output json
	json.NewEncoder(w).Encode(h.c.song)
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL)
}

func publishAPI() {
	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage)
	sh := songsHandler{c: &cache{}}
	rs := statusHandler{c: &cache{}}
	apiRouter.HandleFunc("/api", sh.returnSongs)
	apiRouter.HandleFunc("/api/status/", rs.returnStatus)
	apiRouter.HandleFunc("/post/song", sh.readPost).
		Name("putsong").
		Queries("file", "{file}")
	apiRouter.HandleFunc("/post/status/{mount}", rs.readStatus).
		Name("putstatus").
		Queries("status", "{status}")

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
