package url

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	"mvdan.cc/xurls/v2"
	"os"
	"strings"
	"sync"
)

func skipUrlRuler(url string) bool {
	for _, rule := range config.SkipUrlRulers {
		if strings.Contains(url, rule) {
			return true
		}
	}
	return false
}

func getFileUrl(f *os.File) ([]string, error) {
	// only read one appear(https://) file
	buffer := make([]byte, 1024*1024)
	n, err := f.Read(buffer)
	if err != nil || n == 0 {
		return nil, err
	}
	defer func() {
		if fErr := f.Close(); fErr != nil {
			fmt.Printf("getFileUrl Close:%v\n", err)
		}
	}()
	urls := xurls.Relaxed().FindAllString(string(buffer), -1)
	var data []string
	for _, url := range urls {
		if skipUrlRuler(url) {
			continue
		}
		if strings.Contains(url, "https") {
			data = append(data, url)
		}
	}
	return data, nil
}

func removeDuplicateString(inStrs []string) []string {
	tmp := make(map[string]struct{}, 3)
	for _, str := range inStrs {
		tmp[str] = struct{}{}
	}
	var outStrs []string
	for str := range tmp {
		outStrs = append(outStrs, str)
	}
	return outStrs
}

// adapt file is total left.
func GetUrlsInFiles(fNames []string) (map[string][]string, error) {
	fileUrls := make(map[string][]string, 3)

	var wg sync.WaitGroup
	var errs []error
	for _, fn := range fNames {
		osFile, err := os.OpenFile(fn, os.O_RDONLY, 0644)
		if err != nil {
			fmt.Printf("GetUrlsInFiles:%v\n", err)
			return nil, err
		}
		wg.Add(1)
		go func(name string) {
			fUrl, fErr := getFileUrl(osFile)
			if fErr != nil {
				fmt.Printf("GetUrlsInFiles:%v\n", err)
				errs = append(errs, fErr)
				return
			}
			fileUrls[name] = fUrl
			defer wg.Done()
			return
		}(fn)
	}
	wg.Wait()

	if len(errs) != 0 {
		return nil, fmt.Errorf("%v", errs)
	}
	// Remove duplicate url
	for name, urls := range fileUrls {
		newUrls := removeDuplicateString(urls)
		fileUrls[name] = newUrls
	}
	return fileUrls, nil
}

// https://github.com/kubernetes/kubernetes/blob/release-1.25/cmd/genutils/genutils.go
func GetNewContentFromUrl(urls []string, newVersion string) ([]byte, error) {
	if len(urls) != 1 {
		return nil, fmt.Errorf("too many url:%v", urls)
	}

	url := urls[0]
	var newUrl string
	if strings.Contains(url, "kubernetes") {
		spUrls := strings.Split(url, "/")
		for i, sp := range spUrls {
			if strings.Contains(sp, "release-") {
				spUrls[i] = "release-" + config.K8sSynVersion
				break
			}
		}
		newUrl = strings.Join(spUrls, "/")
		fmt.Printf("old==%v\n new==%v\n", url, newUrl)
	}

	return nil, nil
}
