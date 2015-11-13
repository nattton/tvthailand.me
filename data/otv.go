package data

import (
	"encoding/json"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/facebookgo/httpcontrol"
	"github.com/code-mobi/tvthailand.me/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	OtvDomain     = "http://api.otv.co.th/api/index.php/v3"
	OtvDevCode    = "53336900268229151911"
	OtvSecretKey  = "8540c45823b738220ab09764645e0c82"
	OtvAppID      = "75"
	OtvAppVersion = "1.0"
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

type OtvEpisodePlay struct {
	SeasonDetail  OtvSeasonDetail  `json:"season_detail"`
	Detail        string           `json:"detail"`
	Title         string           `json:"name_th"`
	Thumbnail     string           `json:"thumbnail"`
	EpisodeDetail OtvEpisodeDetail `json:"episode_detail"`
}

type OtvSeasonDetail struct {
	ContentSeasonID string `json:"content_season_id"`
	Title           string `json:"name_th"`
}

type OtvEpisodeDetail struct {
	EpisodeID string        `json:"episode_id"`
	Detail    string        `json:"detail"`
	Title     string        `json:"name_th"`
	Thumbnail string        `json:"cover"`
	Date      string        `json:"date"`
	PartItems []OtvPartItem `json:"part_items"`
}

type OtvPartItem struct {
	ID         string `json:"id"`
	Title      string `json:"name_th"`
	IframeHTML string `json:"stream_url"`
	Cover      string `json:"cover"`
	Thumbnail  string `json:"thumbnail"`
}

func GetOTVEpisodelist(contentID string) (responseBody []byte, otvEpisode OtvEpisode, err error) {
	cacheTime := 5 * time.Minute
	keyOTVEpisodelist := fmt.Sprintf("OTV/Episodelist/%s/0/50", contentID)
	redisClient := utils.OpenRedis()
	jsonResult, err := redisClient.Get(keyOTVEpisodelist).Result()
	if err != nil {
		apiURL := fmt.Sprintf("%s/Episodelist/index/%s/%s/%s/%s/%s/%d/%d", OtvDomain, OtvDevCode, OtvSecretKey, OtvAppID, OtvAppVersion, contentID, 0, 50)
		client := &http.Client{
			Transport: &httpcontrol.Transport{
				RequestTimeout: time.Minute,
				MaxTries:       3,
			},
		}
		resp, err := client.Get(apiURL)
		if err != nil {
			log.Println(err)
			return responseBody, otvEpisode, err
		}
		defer resp.Body.Close()
		responseBody, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return responseBody, otvEpisode, err
		}

		errRedis := redisClient.Set(keyOTVEpisodelist, string(responseBody), cacheTime).Err()
		if errRedis != nil {
			log.Println(errRedis)
		}
	} else {
		responseBody = []byte(jsonResult)
	}

	err = json.Unmarshal(responseBody, &otvEpisode)
	if err != nil {
		log.Println(err)
		return
	}
	for index := range otvEpisode.EpisodeList {
		Date, errT := time.Parse(DateFMT, otvEpisode.EpisodeList[index].Date)
		if errT != nil {
			log.Println(errT)
		}
		otvEpisode.EpisodeList[index].Date = Date.Format(DateLongFMT)
	}
	return
}

func GetOTVEpisodePlay(episodeID string, isMobile bool) (responseBody []byte, otvEpisodePlay OtvEpisodePlay, err error) {
	width := "800"
	height := "340"
	if isMobile {
		width = "320"
		height = "200"
	}
	apiURL := fmt.Sprintf("%s/Episode/oplay", OtvDomain)
	formVal := url.Values{
		"dev_code":    {OtvDevCode},
		"dev_key":     {OtvSecretKey},
		"app_id":      {OtvAppID},
		"app_version": {OtvAppVersion},
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
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	responseBody, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(responseBody, &otvEpisodePlay)
	if err != nil {
		log.Println(err)
		return
	}
	return
}
