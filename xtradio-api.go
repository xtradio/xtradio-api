package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
)

type Song struct {
    Title string `json:"song"`
    Artist string `json:"artist"`
    Show string `json:"show"`
    Image string `json:"image"`
    Album string `json:"album"`
    Length string `json:"length"`
    Secs string `json:"secs"`
    Remaining string `json:"remaining"`
    Share string `json:"share"`
    Url string `json:"url"`
    }

type Songs []Song
    
func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "XTRadio API.")
    fmt.Println("Endpoint Hit: homePage")
}

func readPost(w http.ResponseWriter, r *http.Request){
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

    fmt.Fprintf(w, "Song: %v", vars["file"])
    fmt.Println("readPost hit:", vars["file"], song.Artist, song.Title)
}

func returnSongs(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnApi")
    w.Header().Set("Access-Control-Allow-Origin", "*") 
    var song Song
    songs := Songs{
// Doesn't work:        
    Song{Title: song.Title,
         Artist: "Rihanna & Drake",
         Show: "XTRadio" ,
         Image: "https://img.xtradio.org/tracks/2809999333.jpg",
         Album: "",
         Length: "04:01",
         Secs: "241", 
         Remaining: "217",
         Share: "https://soundcloud.com/sylow-75372780/rihanna-drake-work-sylow-remix-cover-by-reynolds-heesters",
         Url: "Rihanna--Drake-Work-Sylow-Remix-feat-Reynolds--Heesters"},
    }    
    json.NewEncoder(w).Encode(songs)
}

func publishApi() {
    apiRouter := mux.NewRouter().StrictSlash(true)
    apiRouter.HandleFunc("/", homePage)
    apiRouter.HandleFunc("/api", returnSongs)
    apiRouter.HandleFunc("/post/song", readPost).
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
