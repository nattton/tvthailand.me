package data

import (
	"encoding/json"
	"fmt"
	"github.com/facebookgo/httpcontrol"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	OTV_DOMAIN      = "http://api.otv.co.th/api/index.php/v3"
	OTV_DEVCODE     = "53336900268229151911"
	OTV_SECRETKEY   = "8540c45823b738220ab09764645e0c82"
	OTV_APP_ID      = "75"
	OTV_APP_VERSION = "1.0"
)

type OtvEpisode struct {
	ContentSeasonID string           `json:"content_season_id"`
	NameTh          string           `json:"name_th"`
	Detail          string           `json:"detail"`
	ModifiedDate    string           `json:"modified_date"`
	Thumbnail       string           `json:"thumbnail"`
	EpisodeList     []OtvEpisodeList `json:"episode_list"`
}

type OtvEpisodeList struct {
	EpisodeID string `json:"episode_id"`
	Detail    string `json:"detail"`
	NameTh    string `json:"name_th"`
	Thumbnail string `json:"thumbnail"`
	Date      string `json:"date"`
}

func GetOTVEpisodelist(contentID string) (otvEpisode OtvEpisode) {
	apiURL := fmt.Sprintf("%s/Episodelist/index/%s/%s/%s/%s/%s/%d/%d", OTV_DOMAIN, OTV_DEVCODE, OTV_SECRETKEY, OTV_APP_ID, OTV_APP_VERSION, contentID, 0, 50)
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.Get(apiURL)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(body, &otvEpisode)
	if err != nil {
		fmt.Println("JSON Parser Error : ", apiURL)
		panic(err)
	}
	return
}

func GetOTVEpisodePlay(episodeID string, isMobile bool) string {
	width := "800"
	height := "460"
	if isMobile {
		width = "320"
		height = "200"
	}
	apiURL := fmt.Sprintf("%s/Episode/oplay", OTV_DOMAIN)
	formVal := url.Values{
		"dev_code":    {OTV_DEVCODE},
		"dev_key":     {OTV_SECRETKEY},
		"app_id":      {OTV_APP_ID},
		"app_version": {OTV_APP_VERSION},
		"ep_id":       {episodeID},
		"width":       {width},
		"height":      {height},
	}
	client := &http.Client{
		Transport: &httpcontrol.Transport{
			RequestTimeout: time.Minute,
			MaxTries:       3,
		},
	}
	resp, err := client.PostForm(apiURL, formVal)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(body)
}
