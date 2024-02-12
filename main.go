package main

import (
	"log"
	"net/http"

	_ "embed"
)

//go:embed static/index.html
var index []byte

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/download/", http.StripPrefix("/download/", fs))
	http.HandleFunc("/", serveIndex)
	log.Fatal(http.ListenAndServe(":8910", nil))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
