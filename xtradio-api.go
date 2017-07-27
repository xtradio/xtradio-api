package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    "sync"
)

type Song struct {
    Title string `json:"song"`
    Artist string `json:"artist"`
    Show string `json:"show"`
    Image string `json:"image"`
    Album string `json:"album"`
    Length string `json:"length"`
    Secs int `json:"secs"`
    Remaining int `json:"remaining"`
    Share string `json:"share"`
    Url string `json:"url"`
    }

type cache struct {
	sync.RWMutex
	song Song
}

type songsHandler struct {
    c *cache
}

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "XTRadio API.")
    fmt.Println("Endpoint Hit: homePage")
}

func (h songsHandler) readPost(w http.ResponseWriter, r *http.Request){
    vars := mux.Vars(r)
    var song Song
    // Open and connect do DB
    db, err := sql.Open("mysql", "root:test@tcp(127.0.0.1:3306)/radio?charset=utf8")
    checkErr(err)

    // Fetch details for the track
    query := db.QueryRow("SELECT artist, title, album, lenght, share, url, image FROM details WHERE filename=?", vars["file"])

    // Save it to the "Songs" struct
    err = query.Scan(&song.Artist, &song.Title, &song.Album, &song.Length, &song.Share, &song.Url, &song.Image)
    checkErr(err)

    db.Close()

    h.c.Lock()
    defer h.c.Unlock()
    h.c.song = song

    fmt.Fprintf(w, "Song: %v", vars["file"])
    fmt.Println("readPost hit:", vars["file"], song.Artist, song.Title)
}

func (h songsHandler) returnSongs(w http.ResponseWriter, r *http.Request){
    h.c.RLock()
    defer h.c.RUnlock()

    fmt.Println("Endpoint Hit: returnApi")
    w.Header().Set("Access-Control-Allow-Origin", "*") 

    json.NewEncoder(w).Encode(h.c.song)
}

func publishApi() {
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
    publishApi()
}
