package main

import (
	"log"
	"net/http"
	"net/url"

	"github.com/prometheus/client_golang/prometheus"
)

func constructTuneinURL(endpoint string, artist string, title string) string {
	partnerid := getEnv("TUNEIN_PARTNER_ID")
	partnerkey := getEnv("TUNEIN_PARTNER_KEY")
	stationid := getEnv("TUNEIN_STATION_ID")
	if partnerid == "" || partnerkey == "" || stationid == "" {
		return ""
	}

	var URL *url.URL
	URL, _ = url.Parse(endpoint)

	parameters := url.Values{}
	parameters.Add("partnerId", partnerid)
	parameters.Add("partnerKey", partnerkey)
	parameters.Add("id", stationid)
	parameters.Add("artist", artist)
	parameters.Add("title", title)
	URL.RawQuery = parameters.Encode()

	return URL.String()
}

func tuneinAPI(artist string, title string) {
	tuneinEndpoint := "http://air.radiotime.com/Playing.ashx"
	callURL := constructTuneinURL(tuneinEndpoint, artist, title)

	if callURL == "" {
		log.Println("No tunein creds, skipping")
		return
	}

	res, err := http.Get(callURL)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode == 200 {
		tuneinSubmission.With(prometheus.Labels{"artist": artist, "title": title}).Inc()
		log.Printf("Successfull TuneIn submission: %s - %s", artist, title)
	} else {
		log.Println("Tunein submission failed, error code: ", res.StatusCode)
	}
}
