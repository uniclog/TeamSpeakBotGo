package vlc

import (
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

const (
	TrackInfoUrl     string = "http://127.0.0.1:8080/requests/status.xml"
	NextStationUrl   string = "http://127.0.0.1:8080/requests/status.xml?command=pl_next"
	PrevStationUrl   string = "http://127.0.0.1:8080/requests/status.xml?command=pl_previous"
	PlayStationIdUrl string = "http://127.0.0.1:8080/requests/status.xml?command=pl_play&id=%s"
)

func NextStation() {
	vlcControl(NextStationUrl)
}

func PrevStation() {
	vlcControl(PrevStationUrl)
}

func PlayStationById(id string) {
	vlcControl(fmt.Sprintf(PlayStationIdUrl, id))
}

func vlcControl(url string) []byte {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(":1")))

	client := &http.Client{}
	resp, err1 := client.Do(req)
	if err1 != nil {
	    return nil
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	return body
}

func GetTrackInfo() (string, string) {
	body := vlcControl(TrackInfoUrl)
	var status StatusResponse
	_ = xml.Unmarshal(body, &status)

	var nowPlaying, title string

	for _, category := range status.Information.Categories {
		if category.Name == "meta" {
			for _, info := range category.Infos {
				if info.Name == "title" {
					title = info.Value
				} else if info.Name == "now_playing" {
					nowPlaying = info.Value
				}
			}
		}
	}
	return title, nowPlaying
}

type StatusResponse struct {
	Information struct {
		Categories []Category `xml:"category"`
	} `xml:"information"`
}

type Category struct {
	Name  string `xml:"name,attr"`
	Infos []Info `xml:"info"`
}

type Info struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}
