package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
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
	}else if r.Method == "POST" {
		//接收文件流
		file, head, err := r.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data, err: %s\n", err.Error())
			return
		}
		fmt.Printf("file: %s\n", head.Filename)
		defer file.Close()

		newFile, err := os.Create("/tmp/" + head.Filename)
		if err != nil {
			fmt.Printf("failed to create file, err:%s\n" , err.Error())
			return
		}
		defer newFile.Close()

		_, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("faild to save data into file, err:%s\n", err.Error())
			return
		}
		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "upload success!")
}
