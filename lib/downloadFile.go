package lib

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

const DownloadsDir = "./temp/"

func DownloadFile(url, dir string, ch chan bool) error {
	defer func() {
		<-ch
	}()
	stringArray := strings.Split(url, "/")
	fileName := stringArray[len(stringArray)-1]

	os.MkdirAll(DownloadsDir+dir, os.ModePerm)

	out, err := os.Create(DownloadsDir + dir + "/" + fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}
	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

// concurrent 下载并发数
func DownloadSeriesFile(series []Series, dir string, concurrent int) {
	// 下载文件
	fmt.Println("downloading...")
	ch := make(chan bool, concurrent)
	for _, serie := range series {
		for _, instance := range serie.InstanceList {
			ch <- true
			go DownloadFile(instance.ImageId, dir, ch)
		}
	}

	// 压缩文件
	fmt.Println("zipping...")
	os.MkdirAll(ZipsDir, os.ModePerm)
	eee := Zipit(DownloadsDir+dir, ZipsDir+dir+".zip")
	fmt.Println(eee)

	// 删除源文件
	fmt.Println("deleting origin files...")
	os.RemoveAll(DownloadsDir + dir)

	fmt.Println("---- FINISHED ----")
}
