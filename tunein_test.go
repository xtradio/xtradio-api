package main

import (
	"net/url"
	"os"
	"testing"
)

func TestConstructTuneinURL(t *testing.T) {
	testURLHost := "http://google.com"
	testArtist := "test"
	testTitle := "title"

	testPartnerID := "foo"
	testPartnerKey := "bar"
	testStationID := "baz"

	os.Setenv("TUNEIN_PARTNER_ID", testPartnerID)
	os.Setenv("TUNEIN_PARTNER_KEY", testPartnerKey)
	os.Setenv("TUNEIN_STATION_ID", testStationID)

	var testURL *url.URL

	testURL, _ = url.Parse(testURLHost)

	testParams := url.Values{}
	testParams.Add("partnerId", testPartnerID)
	testParams.Add("partnerKey", testPartnerKey)
	testParams.Add("id", testStationID)
	testParams.Add("artist", testArtist)
	testParams.Add("title", testTitle)

	testURL.RawQuery = testParams.Encode()

	getURL := constructTuneinURL(testURLHost, testArtist, testTitle)

	// fmt.Println(getURL.String())

	if getURL != testURL.String() {
		t.Errorf("URL Construct test: got %s, expected %s", getURL, testURL.String())
	}

	os.Unsetenv("TUNEIN_PARTNER_ID")
	os.Unsetenv("TUNEIN_PARTNER_KEY")
	os.Unsetenv("TUNEIN_STATION_ID")

}

func TestConstructTuneinURLNoEnvVars(t *testing.T) {
	getData := constructTuneinURL("url.com", "test", "test")

	if getData != "" {
		t.Errorf("RUL Construct no ENV Key: expected empty string, got %s", getData)
	}
}
