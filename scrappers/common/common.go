package common

import (
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func HitGetRequest(path string) []byte {
	client := http.Client{Timeout: time.Second * 10}
	resp, err := client.Get(path)
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close() // nolint: errcheck
	byteValue, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err.Error())
		return nil
	}
	return byteValue
}
