package main

import (
	"fmt"
	"github.com/chiefsh/goPan/handler"
	"net/http"
)

func main() {
	http.HandleFunc("/file/upload", handler.UploadHandler)
	err := http.ListenAndServe(":8080", nil)
	fmt.Println("server start at: localhost:8080")
	if err != nil {
		fmt.Printf("Failed to start server on localhost:8080 ...")
	}
}
