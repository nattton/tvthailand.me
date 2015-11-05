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
	ThumbnailURLCategory = "http://thumbnail.instardara.com/category/"
	ThumbnailURLChannel  = "http://thumbnail.instardara.com/channel/"
	ThumbnailURLRadio    = "http://thumbnail.instardara.com/radio/"
	ThumbnailURLTv       = "http://thumbnail.instardara.com/tv/"
	ThumbnailURLPoster   = "http://thumbnail.instardara.com/poster/"
)

const SECRET_SALT = "Cod3M0b!"

var throttle = make(chan int, MaxConcurrency)

func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(SECRET_SALT+plaintext)))
	return
}
