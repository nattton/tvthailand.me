package analytic

import (
	"encoding/json"
	"fmt"
	"github.com/code-mobi/tvthailand.me/Godeps/_workspace/src/github.com/jinzhu/gorm"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"sync"
)

const MaxConcurrency = 14

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
}

type Result struct {
	Title       string
	RowAffected int64
}

type Show struct {
	ID        int `gorm:"primary_key"`
	Title     string
	ViewCount int
}

func UpdateView(db *gorm.DB, jsonByte []byte) {
	shows := getShow(jsonByte)
	if len(shows) > 0 {
		resetView(db)
	}

	var wg sync.WaitGroup
	for _, show := range shows {
		throttle <- 1
		wg.Add(1)
		go updateViewCount(db, show, &wg, throttle)
	}
	wg.Wait()
}

func resetView(db *gorm.DB) {
	err := db.Model(Show{}).UpdateColumn("view_count", 0).Error
	if err != nil {
		log.Fatal(err)
	}
}

func getShow(b []byte) []*ShowItem {
	var analytics Analytics
	err := json.Unmarshal(b, &analytics)
	if err != nil {
		panic(err)
	}
	var shows []*ShowItem
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

		shows = append(shows, &ShowItem{displayKey, dataValue})
	}
	return shows
}

func updateViewCount(db *gorm.DB, showItem *ShowItem, wg *sync.WaitGroup, throttle chan int) {
	defer wg.Done()
	err := db.Model(Show{}).Where("title = ?", showItem.Title).UpdateColumn("view_count", showItem.ViewCount).Error
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(showItem.Title, showItem.ViewCount)
	<-throttle
}
