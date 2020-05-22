package main

import (
	"fmt"
	"regexp"
	"strconv"
)

// UpcomingSongs details
type UpcomingSongs struct {
	Title  string  `json:"song"`
	Artist string  `json:"artist"`
	Image  string  `json:"image"`
	Length float64 `json:"length"`
	Share  string  `json:"share"`
}

func fetchUpcomingSongsFromDb(list []string) []UpcomingSongs {

	var r []UpcomingSongs

	for _, v := range list {

		var u UpcomingSongs

		re := regexp.MustCompile(`artist="(.*?)",`)
		u.Artist = re.FindStringSubmatch(v)[1]

		re = regexp.MustCompile(`title="(.*?)",`)
		u.Title = re.FindStringSubmatch(v)[1]

		re = regexp.MustCompile(`share="(.*?)"`)
		u.Share = re.FindStringSubmatch(v)[1]

		re = regexp.MustCompile(`length="(.*?)",`)
		u.Length, _ = strconv.ParseFloat(re.FindStringSubmatch(v)[1], 10)

		re = regexp.MustCompile(`image="(.*?)",`)
		u.Image = re.FindStringSubmatch(v)[1]
		if u.Image == "" {
			u.Image = "default.png"
		}

		u.Image = fmt.Sprintf("https://img.xtradio.org/tracks/%s", u.Image)

		r = append(r, u)

	}

	return r

}

func upcomingSongs() ([]UpcomingSongs, error) {
	var dbParsedData []UpcomingSongs
	command := "playlist(dot)txt.next"
	data, err := telnet(command)

	if err != nil {
		return dbParsedData, err
	}
	dbParsedData = fetchUpcomingSongsFromDb(data)

	return dbParsedData, nil
}
