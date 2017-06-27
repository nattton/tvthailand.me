package data

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
)

const (
	DateLongFMT    = "02 Jan 2006"
	DateFMT        = "2006-01-02"
	DateTimeFMT    = "2006-01-02 15:04:05"
	LimitRow       = 40
	MaxConcurrency = 8
)

const (
	ThumbnailURLCategory = "http://thumbnail.instardara.com/category/"
	ThumbnailURLChannel  = "http://thumbnail.instardara.com/channel/"
	ThumbnailURLRadio    = "http://thumbnail.instardara.com/radio/"
	ThumbnailURLTv       = "http://thumbnail.instardara.com/tv/"
	ThumbnailURLPoster   = "http://thumbnail.instardara.com/poster/"
	YoutubeViewURL       = "https://www.youtube.com/watch?v="
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

var re = strings.NewReplacer(
	"+", "-",
	"=", ",",
	"a", "!",
	"b", "@",
	"c", "#",
	"d", "$",
	"e", "%",
	"f", "^",
	"g", "&",
	"h", "*",
	"i", "(",
	"j", ")",
	"k", "{",
	"l", "}",
	"m", "[",
	"n", "]",
	"o", ":",
	"p", ";",
	"q", "<",
	"r", ">",
	"s", "?",
)

func EncryptVideo(videoID string) string {
	encrypt := base64.StdEncoding.EncodeToString([]byte(videoID))
	return re.Replace(encrypt)
}
