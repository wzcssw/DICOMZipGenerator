package main

import (
	"DICOMZipGenerator/lib"
	"encoding/json"
	"runtime"

	"github.com/gin-gonic/gin"
)

// Next:
//	 加入Redis

// http://localhost:8080/api/dicom/meeting138
// const DICOMServerURL string = "http://47.93.132.62/api/getByFilmNo" // test
const DICOMServerURL string = "http://dicomup.tongxinyiliao.com/api/getByFilmNo"

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	r := gin.Default()

	r.GET("/api/dicom/:filmno", func(c *gin.Context) {
		m := make(map[string]string)
		m["filmno"] = c.Param("filmno")
		responsedDataBody := lib.SendDicomAPIRequest(DICOMServerURL, m)

		var responsedData lib.DicomAPIRequest
		json.Unmarshal([]byte(responsedDataBody), &responsedData)
		lib.DownloadSeriesFile(responsedData.List, c.Param("filmno"), 30)

		c.JSON(200, gin.H{
			"data":    responsedDataBody,
			"msg":     "OK",
			"success": true,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

// 一个DICOM文件
// https://dicom.tongxinyiliao.com/2018/06/06/8/a47400ce7333880a4073022777d84b99.dcm
