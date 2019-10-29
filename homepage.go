package main

import (
	"fmt"
	"net/http"
	"time"
)

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "XTRadio API.")
	fmt.Println(time.Now(), r.RemoteAddr, r.Method, r.URL)
}
