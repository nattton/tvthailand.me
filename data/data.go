package data

import (
	"crypto/sha1"
	"fmt"
)

const (
	DateLongFMT    = "02 Jan 2006"
	DateFMT        = "2006-01-02"
	DateTimeFMT    = "2006-01-02 15:04:05"
	LimitRow       = 40
	MaxConcurrency = 8
)

const (
	ThumbnailURLCategory = "https://thumbnail.tvthailand.me/category/"
	ThumbnailURLChannel  = "https://thumbnail.tvthailand.me/channel/"
	ThumbnailURLRadio    = "https://thumbnail.tvthailand.me/radio/"
	ThumbnailURLTv       = "https://thumbnail.tvthailand.me/tv/"
	ThumbnailURLPoster   = "https://thumbnail.tvthailand.me/poster/"
)

const SECRET_SALT = "Cod3M0b!"

var throttle = make(chan int, MaxConcurrency)

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(SECRET_SALT+plaintext)))
	return
}

type PageInfo struct {
	PreviousPage   int32 `json:"previousPage,omitempty"`
	NextPage       int32 `json:"nextPage,omitempty"`
	TotalResults   int32 `json:"totalResults"`
	ResultsPerPage int32 `json:"resultsPerPage"`
}
