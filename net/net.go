package net

import (
	"fmt"
	"io/ioutil"
	"net/http"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func getURLContent(newUrl string) ([]byte, error) {
	resp, err := http.Get(newUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("resp is: %v", resp.Status)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GetURLsContent(newUrls []string) ([]byte, error) {
	var value []byte
	var parseErrs []error
	for _, url := range newUrls {
		if len(newUrls) > 1 {
			value = append(value, []byte(url)...)
			value = append(value, '\n')
		}
		data, err := getURLContent(url)
		value = append(value, data...)
		parseErrs = append(parseErrs, err)
	}
	return value, utilerrors.NewAggregate(parseErrs)
}
