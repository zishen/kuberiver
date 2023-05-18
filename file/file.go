package file

import (
	"github.com/zishen/kuberiver/config"
	hwlog "github.com/zishen/kuberiver/log"
	"github.com/zishen/kuberiver/url"
	"io/ioutil"
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
			hwlog.RunLog.Errorf("walking bootstrap dirs failed: %v: %v", path, err)
			return err
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

func creatNewDataFile(fName string, data []byte, changeFlag bool) error {
	// get abs path.
	path, pathErr := filepath.Abs(fName)
	if pathErr != nil {
		return pathErr
	}
	index := strings.LastIndex(path, string(os.PathSeparator))
	dir := path[:index]
	// create multilevel directories
	dirErr := os.MkdirAll(dir, os.ModePerm)
	if dirErr != nil {
		return dirErr
	}

	newFileName := fName
	if !changeFlag {
		newFileName = strings.Replace(fName, ".go", ".html", 1)
	}

	out, err := os.Create(newFileName)
	if err != nil {
		return err
	}
	hwlog.RunLog.Debugf("creat new file:%v success.", newFileName)
	defer out.Close()
	return ioutil.WriteFile(newFileName, data, 0644)
}

func SetNewFilesByURLs(prefix, newPath string, fileUrls map[string][]string) error {
	for oldName, urls := range fileUrls {
		hwlog.RunLog.Debugf("begin SetNewFilesByURLs:%v", oldName)
		if len(urls) == 0 {
			hwlog.RunLog.Errorf("SetNewFilesByURLs %v url is null", oldName)
			continue
		}
		newName := strings.Replace(oldName, prefix, newPath, 1)

		data, urlErr, changeFlag := url.GetNewContentFromUrl(urls, config.K8sSynVersion)
		if urlErr != nil {
			hwlog.RunLog.Errorf("GetNewContentFromUrl:%v", urlErr)
			continue
		}
		if createErr := creatNewDataFile(newName, data, changeFlag); createErr != nil {
			hwlog.RunLog.Errorf("creatNewDataFile:%v", createErr)
			continue
		}
		hwlog.RunLog.Debugf("end SetNewFilesByURLs:%v", oldName)
	}
	return nil
}
