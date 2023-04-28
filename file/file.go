package file

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	"github.com/zishen/kuberiver/url"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var ignorePrefixes = []string{
	".",
	"_",
	"#",
}

func GetAllFilesPath(oldPath string) ([]string, error) {
	var files []string
	walkErr := filepath.Walk(oldPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatalf("walking bootstrap dirs failed: %v: %v", path, err)
		}
		// skip dir file
		if info.IsDir() {
			return nil
		}
		// skip base system file
		name := filepath.Base(path)
		for _, pre := range ignorePrefixes {
			if strings.HasPrefix(name, pre) {
				return nil
			}
		}
		// skip test file.
		if strings.Contains(path, "_test.go") {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, walkErr
}

func creatNewDataFile(fName string, data []byte) error {
	fmt.Printf("creatNewDataFile:%v\n", fName)
	return nil
}

func SetNewFilesByURLs(prefix, newPath string, fileUrls map[string][]string) error {
	for oldName, urls := range fileUrls {
		newName := strings.Replace(oldName, prefix, newPath, 1)

		data, urlErr := url.GetNewContentFromUrl(urls, config.K8sSynVersion)
		if urlErr != nil {
			fmt.Printf("SetNewFilesByURLs:%v\n", urlErr)
			continue
		}
		if createErr := creatNewDataFile(newName, data); createErr != nil {
			return urlErr
		}
	}
	return nil
}
