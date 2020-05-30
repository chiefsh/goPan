package handler

import (
	"io"
	"io/ioutil"
	"net/http"
)

// 处理文件长传
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method=="GET" {
		// 返回上传的html
		data, err := ioutil.ReadFile("./static/view/index.html")
		if err != nil {
			io.WriteString(w, "internel server error")
			return
		}
		io.WriteString(w, string(data))
	}else {
		//接收文件流
	}
}
