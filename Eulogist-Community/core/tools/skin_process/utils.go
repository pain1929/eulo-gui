package skin_process

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

// ...
func IsZIPFile(fileData []byte) bool {
	return len(fileData) >= 4 && bytes.Equal(fileData[0:4], []byte("PK\x03\x04"))
}

// 从 url 指定的网址下载文件，
// 并返回该文件的二进制形式
func DownloadFile(url string) (result []byte, err error) {
	// 获取 HTTP 响应
	httpResponse, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("DownloadFile: %v", err)
	}
	defer httpResponse.Body.Close()
	// 读取文件数据
	result, err = io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, fmt.Errorf("DownloadFile: %v", err)
	}
	// 返回值
	return
}
