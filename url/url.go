package url

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	hwlog "github.com/zishen/kuberiver/log"
	"github.com/zishen/kuberiver/net"
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
			hwlog.RunLog.Errorf("getFileUrl Close:%v", err)
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
			hwlog.RunLog.Errorf("OpenFile:%v", err)
			return nil, err
		}
		wg.Add(1)
		go func(name string) {
			fUrl, fErr := getFileUrl(osFile)
			if fErr != nil {
				hwlog.RunLog.Errorf("getFileUrl:%v", err)
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

// https://github.com/kubernetes/kubernetes/blob/release-1.26/pkg/scheduler/framework/types.go
// https://raw.githubusercontent.com/kubernetes/kubernetes/release-1.26/pkg/scheduler/framework/types.go
func getNewUrl(url string, newVersion string) (string, bool) {
	newUrl := url
	hwlog.RunLog.Debugf("old url [%v]", newUrl)
	if strings.Contains(url, "github.com") {
		spUrls := strings.Split(url, "/")
		for i, sp := range spUrls {
			if strings.Contains(sp, "github.com") {
				spUrls[i] = config.GitHubPreRawURL
				continue
			}
			if strings.Contains(sp, "release-") {
				spUrls[i] = "release-" + newVersion
				continue
			}
		}
		for i, sp := range spUrls {
			if strings.Contains(sp, "blob") {
				spUrls = append(spUrls[:i], spUrls[i+1:]...) // delete
				continue
			}
		}
		newUrl = strings.Join(spUrls, "/")
		hwlog.RunLog.Debugf("new url [%v]", newUrl)
	} else {
		hwlog.RunLog.Infof("no change url [%v]", newUrl)
		return newUrl, false
	}
	return newUrl, true
}

// https://github.com/kubernetes/kubernetes/blob/release-1.25/cmd/genutils/genutils.go
func GetNewContentFromUrl(urls []string, newVersion string) ([]byte, error, bool) {
	if len(urls) != 1 {
		return nil, fmt.Errorf("too many url:%v", urls), false
	}

	newUrl, changeFlag := getNewUrl(urls[0], newVersion)

	data, err := net.GetURLContent(newUrl)
	return data, err, changeFlag
}
