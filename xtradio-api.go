package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XTRadio API.")
	log.Println(r.RemoteAddr, r.Method, r.URL)
}

func publishAPI() {
	apiRouter := mux.NewRouter().StrictSlash(true)
	apiRouter.HandleFunc("/", homePage)
	apiRouter.HandleFunc("/search", songSearch)
	apiRouter.HandleFunc("/v1/song/list", songList).
		Methods("GET")
	apiRouter.HandleFunc("/v1/song/upload", songUpload).
		Methods("POST")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	log.Fatal(http.ListenAndServe(":10000", c.Handler(apiRouter)))
}

func main() {
	log.Println("Rest API v2.0 - Mux Routers")
	publishAPI()
}
