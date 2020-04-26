package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2/clientcredentials"

	"github.com/yanatan16/golang-soundcloud/soundcloud"
	"github.com/zmb3/spotify"
)

// SongResults holds the data returned from spotify and soundcloud
type SongResults struct {
	// Biggest image of the song from the album art
	ImageURL string `json:"img"`
	// Name of the artist of the song
	Artist string `json:"artist"`
	// Name of the title of the song
	Title string `json:"title"`
	// Total duration in seconds of the song
	Duration uint64 `json:"duration"`
	// The public URL of the song for link back
	SourceURL string `json:"sourceurl"`
	// Name of the user who uploaded the song
	User string `json:"user"`
	// Name of the service where the song information came from (spotify|soundcloud)
	Service string `json:"service"`
	// BPM of the track
	BPM float64 `json:"bpm"`
	// The genre of the song
	Genre string `json:"genre"`
	// License of the song
	License string `json:"license"`
}

func songSearch(w http.ResponseWriter, r *http.Request) {

	var data []SongResults

	fmt.Println(r.FormValue("artist"), r.FormValue("title"))

	artist := r.FormValue("artist")
	title := r.FormValue("title")

	spotify := spotifySearch(artist, title)

	for _, track := range spotify {
		data = append(data, track)
	}

	soundcloud := scSearch(artist, title)

	for _, track := range soundcloud {
		data = append(data, track)
	}

	json.NewEncoder(w).Encode(data)
	log.Println(r.RemoteAddr, r.Method, r.URL)

}

func spotifySearch(artist string, title string) []SongResults {
	var clientID string
	var clientSecret string

	var data SongResults
	var list []SongResults

	clientID, _ = getEnv("SPOTIFY_CLIENT_ID")
	clientSecret, _ = getEnv("SPOTIFY_CLIENT_SECRET")

	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     spotify.TokenURL,
	}

	token, err := config.Token(context.Background())
	if err != nil {
		log.Fatalf("couldn't get token: %v", err)
	}

	search := fmt.Sprintf("%s - %s", artist, title)

	client := spotify.Authenticator{}.NewClient(token)

	results, err := client.Search(search, spotify.SearchTypeTrack)
	if err != nil {
		log.Fatal(err)
	}

	if results.Tracks == nil {
		return list
	}

	for _, item := range results.Tracks.Tracks {

		for _, v := range item.SimpleTrack.Artists {
			data.Artist = v.Name
		}
		data.Title = item.SimpleTrack.Name
		data.Duration = uint64(item.SimpleTrack.Duration)
		album, _ := client.GetAlbum(item.Album.ID)

		for _, a := range album.Artists {
			data.User = a.Name
			continue
		}

		for _, g := range album.Genres {
			data.Genre = g
			continue
		}

		for _, i := range item.Album.Images {
			if i.Height == 640 {
				data.ImageURL = i.URL
			}
		}

		for _, u := range item.ExternalURLs {
			data.SourceURL = u
			continue
		}

		// var audioFeature *spotify.AudioAnalysis
		audioFeature, err := client.GetAudioFeatures(item.SimpleTrack.ID)

		if err != nil {
			fmt.Println("Error getting the analysis of track ", item.SimpleTrack.ID, err)
		}

		for _, f := range audioFeature {
			bpm := math.Round(float64(f.Tempo))
			data.BPM = bpm
			continue
		}

		data.Service = "spotify"

		list = append(list, data)
	}

	return list

}

func scSearch(artist string, title string) []SongResults {

	var data SongResults
	var list []SongResults

	var clientID string
	clientID, _ = getEnv("SOUNDCLOUD_CLIENT_ID")
	api := &soundcloud.Api{
		ClientId: clientID,
	}

	search := fmt.Sprintf("%s - %s", artist, title)

	ret, err := api.Tracks(url.Values{"q": []string{search}})

	if err != nil {
		fmt.Println("SC Error: ", err)
	}

	for _, v := range ret {
		split := strings.Split(v.Title, " - ")
		if len(split) > 1 {
			data.Artist = split[0]
			data.Title = split[1]
		} else {
			data.Title = fmt.Sprintf(v.Title)
		}

		data.Duration = v.Duration

		imageURL := strings.Replace(fmt.Sprintf(v.ArtworkUrl), "large", "t500x500", 1)

		data.ImageURL = imageURL
		data.SourceURL = v.PermalinkUrl
		data.Service = "soundcloud"
		data.Genre = v.Genre
		data.User = v.User.Username
		data.License = v.License
		data.BPM = v.Bpm

		list = append(list, data)

	}

	return list

}
