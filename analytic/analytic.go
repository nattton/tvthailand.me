package analytic

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/data"
	"log"
	"strconv"
	"strings"
	"sync"
)

const MaxConcurrency = 4

var throttle = make(chan int, MaxConcurrency)

type ShowItem struct {
	Title     string
	ViewCount int
	Updated   bool
}

func getShowCSV(b []byte) (shows []ShowItem, err error) {
	r := csv.NewReader(bytes.NewReader(b))
	r.Comma = ','
	r.Comment = '#'
	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for _, row := range records[1:] {
		viewCount, _ := strconv.Atoi(strings.Replace(row[1], ",", "", -1))
		if viewCount > 0 {
			show := ShowItem{
				Title:     row[0],
				ViewCount: viewCount,
			}
			shows = append(shows, show)
			fmt.Println(show.Title, show.ViewCount)
		}
	}
	return
}

func UpdateView(db *gorm.DB, jsonByte []byte) []ShowItem {
	shows, _ := getShowCSV(jsonByte)
	if len(shows) > 0 {
		data.ResetShowViewCount(db)
	}

	var wg sync.WaitGroup
	for index := range shows {
		throttle <- 1
		wg.Add(1)
		go updateViewCount(db, &shows[index], &wg, throttle)
	}
	wg.Wait()
	return shows
}

func updateViewCount(db *gorm.DB, showItem *ShowItem, wg *sync.WaitGroup, throttle chan int) {
	defer wg.Done()
	rowsAffected := data.UpdateShowViewCount(db, showItem.Title, showItem.ViewCount)
	showItem.Updated = rowsAffected > 0
	<-throttle
}
