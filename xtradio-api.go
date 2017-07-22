package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"

)

func homePage(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "XTRadio API.")
    fmt.Println("Endpoint Hit: homePage")
}

func handleRequests() {
    myRouter := mux.NewRouter().StrictSlash(true)
    myRouter.HandleFunc("/", homePage)
    myRouter.HandleFunc("/api", returnSongs)
    log.Fatal(http.ListenAndServe(":10000", myRouter))
}

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
    
func returnSongs(w http.ResponseWriter, r *http.Request){
    songs := Songs{
        Song{Title: "Work (Sylow Remix feat. Reynolds & Heesters)",
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
    fmt.Println("Endpoint Hit: returnApi")
    w.Header().Set("Access-Control-Allow-Origin", "*") 
    json.NewEncoder(w).Encode(songs)
}

func main() {
    fmt.Println("Rest API v2.0 - Mux Routers")
    handleRequests()
}
