package net

import (
	"io/ioutil"
	"net/http"
)

func GetURLContent(newUrl string) ([]byte, error) {
	resp, err := http.Get(newUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}
