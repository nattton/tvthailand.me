package analytic

import (
	"encoding/json"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	"github.com/code-mobi/tvthailand.me/data"
	"strconv"
	"strings"
	"sync"
)

const MaxConcurrency = 4

var throttle = make(chan int, MaxConcurrency)

type Analytics struct {
	Components []*Components `json:"components"`
}

type Components struct {
	DataTable DataTable `json:"dataTable"`
}

type DataTable struct {
	RowClusters []*RowCluster `json:"rowCluster"`
}

type RowCluster struct {
	RowKey []*RowKey `json:"rowKey"`
	Row    []*Row    `json:"row"`
}

type RowKey struct {
	DisplayKey string `json:"displayKey"`
}

type Row struct {
	RowValue []*RowValue `json:"rowValue"`
}

type RowValue struct {
	DataValue string `json:"dataValue"`
}

type ShowItem struct {
	Title     string
	ViewCount int
	Updated   bool
}

type Result struct {
	Title       string
	RowAffected int64
}

func getShow(b []byte) []ShowItem {
	var analytics Analytics
	err := json.Unmarshal(b, &analytics)
	if err != nil {
		panic(err)
	}
	var shows []ShowItem
	rowClusters := analytics.Components[0].DataTable.RowClusters
	if len(rowClusters) == 0 {
		rowClusters = analytics.Components[1].DataTable.RowClusters
	}
	for _, rowCluster := range rowClusters {
		displayKey := rowCluster.RowKey[0].DisplayKey
		dataValue, err := strconv.Atoi(strings.Replace(rowCluster.Row[0].RowValue[0].DataValue, ",", "", -1))
		if err != nil {
			panic(err)
		}

		shows = append(shows, ShowItem{displayKey, dataValue, false})
	}
	return shows
}

func UpdateView(db *gorm.DB, jsonByte []byte) []ShowItem {
	shows := getShow(jsonByte)
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
