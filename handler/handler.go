package handler

import (
	"encoding/json"
	"fmt"
	"github.com/chiefsh/goPan/util"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"github.com/chiefsh/goPan/meta"
	"time"
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

		fileMeta := meta.FileMeta{
			FileName: head.Filename,
			Location: "./tmp/" + head.Filename,
			UploadAt: time.Now().Format("2006-01-02 15:04:05"),
		}
		newFile, err := os.Create(fileMeta.FileName)
		if err != nil {
			fmt.Printf("failed to create file, err:%s\n" , err.Error())
			return
		}
		defer newFile.Close()

		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("faild to save data into file, err:%s\n", err.Error())
			return
		}
		newFile.Seek(0,0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		meta.UpdateFileMeta(fileMeta)

		http.Redirect(w, r, "/file/upload/suc", http.StatusFound)
	}
}

func UploadSucHandler(w http.ResponseWriter, r *http.Request)  {
	io.WriteString(w, "upload success!")
}

func GetFileMetaHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	filehash := r.Form["filehash"][0]
	fmeta := meta.GetFileMeta(filehash)
	data, err := json.Marshal(fmeta)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

func DownloadHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	fsha1 := r.Form.Get("filehash")
	fm := meta.GetFileMeta(fsha1)

	f, err := os.Open(fm.Location)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octect-stream")
	w.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	w.Write(data)
}

func FileMetaUpdateHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	opType := r.Form.Get("op")
	filesha1 := r.Form.Get("filehash")
	newFileName := r.Form.Get("filename")

	if opType != "0" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := meta.GetFileMeta(filesha1)
	curFileMeta.FileName = newFileName
	meta.UpdateFileMeta(curFileMeta)
	w.WriteHeader(http.StatusOK)
}

func FileDeleteHandler(w http.ResponseWriter, r *http.Request)  {
	r.ParseForm()

	filesha1 := r.Form.Get("filehash")
	fMeta := meta.GetFileMeta(filesha1)
	os.Remove(fMeta.Location)
	meta.RemoveFileMeta(filesha1)

	w.WriteHeader(http.StatusOK)
}