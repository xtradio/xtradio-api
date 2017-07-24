package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "github.com/gorilla/mux"
    _ "github.com/go-sql-driver/mysql"
    "database/sql"
    boltDB "github.com/boltdb/bolt"
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

    saveData("artist", song.Artist)
    saveData("title", song.Title)
    saveData("album", song.Album)
    saveData("length", song.Length)
    saveData("share", song.Share)
    saveData("url", song.Url)
    saveData("image", song.Image)

    db.Close()

    fmt.Fprintf(w, "Song: %v", vars["file"])
    fmt.Println("readPost hit:", vars["file"], song.Artist, song.Title)
}

func returnSongs(w http.ResponseWriter, r *http.Request){
    fmt.Println("Endpoint Hit: returnApi")
    w.Header().Set("Access-Control-Allow-Origin", "*") 

    bolt, err := boltDB.Open("my.db", 0600, nil)
    checkErr(err)

    var song []byte

    bolt.View(func(tx *boltDB.Tx) error {
	b := tx.Bucket([]byte("songlist"))
	b.ForEach(func(k, v []byte) error {
		fmt.Printf("key=%s, value=%s\n", k, v)
                return nil
	})
	v := b.Get([]byte("artist"))
	fmt.Printf("The answer is: %s\n", v)
	return nil
    })
    defer bolt.Close()
    json.NewEncoder(w).Encode(song)
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

func saveData(key, value  string) error {
    bolt, err := boltDB.Open("my.db", 0600, nil)
    checkErr(err)

    bolt.Update(func(tx *boltDB.Tx) error {
	b, err := tx.CreateBucketIfNotExists([]byte("songlist"))
        checkErr(err)

	err2 := b.Put([]byte(key), []byte(value))
	if err2 != nil {
		return fmt.Errorf("create value: %s", err2)
	}
        fmt.Println("Data save: ", key, value)
	return nil
    })

    defer bolt.Close()

    return nil
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
