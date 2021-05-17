package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HitGetRequest(path string) ([]byte, int) {
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(path)
	if err != nil {
		log.Println(err)
		return nil, 0
	}
	defer resp.Body.Close() // nolint: errcheck
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil, 0
	}
	return byteValue, resp.StatusCode
}
