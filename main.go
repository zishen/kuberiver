package main

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	"github.com/zishen/kuberiver/file"
	"github.com/zishen/kuberiver/log"
	"github.com/zishen/kuberiver/url"
)

func main() {
	if err := zklog.InitZKLogger(); err != nil {
		fmt.Printf("get InitZKLogger failed: %v", err)
		return
	}
	/*oldPath :="/usr/workspace/go/src/karmada-io/karmada/pkg/util/lifted"*/
	files, fileErr := file.GetAllFilesPath(config.LiftedPath)
	if fileErr != nil {
		zklog.RunLog.Errorf("GetAllFilesPath:%v", fileErr)
		return
	}
	zklog.RunLog.Debugf("GetAllFilesPath:%v", files)

	fileUrls, urlErr := url.GetUrlsInFiles(files)
	if urlErr != nil {
		zklog.RunLog.Errorf("GetUrlsInFiles:%v", urlErr)
		return
	}
	zklog.RunLog.Debugf("get fileUrls::%v", fileUrls)

	conErr := file.SetNewFilesByURLs(config.LiftedPath, config.TmpPath, fileUrls)
	if conErr != nil {
		zklog.RunLog.Errorf("SetNewFilesByURLs:%v", conErr)
		return
	}
}
