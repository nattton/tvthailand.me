package validate

import (
	"log"
	"net/http"
)

func URL(url string) (bool, error) {
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return false, err
	}
	if resp.StatusCode == 404 {
		return false, nil
	}
	return true, nil
}
