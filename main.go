package main

import (
	"fmt"
	"github.com/chiefsh/goPan/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	http.HandleFunc("/file/upload/suc", handler.UploadSucHandler)
	fmt.Println("server start at: localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server on localhost:8080 ...")
	}
}
