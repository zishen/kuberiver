package main

import (
	"fmt"
	"github.com/zishen/kuberiver/config"
	"github.com/zishen/kuberiver/file"
	"github.com/zishen/kuberiver/url"
)

func main() {
	/*oldPath :="/usr/workspace/go/src/karmada-io/karmada/pkg/util/lifted"*/
	files, fileErr := file.GetAllFilesPath(config.LiftedPath)
	if fileErr != nil {
		fmt.Printf("GetAllFilesPath:%v\n", fileErr)
		return
	}
	//fmt.Printf("GetAllFilesPath:%v\n",files)

	fileUrls, urlErr := url.GetUrlsInFiles(files)
	if urlErr != nil {
		fmt.Printf("GetUrlsInFiles:%v\n", urlErr)
		return
	}
	/*	fmt.Println("get fileUrls:")
		for fName, urls := range fileUrls {
			fmt.Printf("========%s========\n%+v\n\n", fName, urls)
		}*/

	conErr := file.SetNewFilesByURLs(config.LiftedPath, config.TmpPath, fileUrls)
	if conErr != nil {
		fmt.Printf("SetNewFilesByURLs:%v\n", conErr)
		return
	}
}
