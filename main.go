package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"

	_ "embed"
)

//go:embed static/index.html
var index []byte

var rootDirectory = "."

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/api/download/", http.StripPrefix("/download/", fs))
	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/list/{path...}", listContent)
	log.Fatal(http.ListenAndServe(":8910", nil))
}

func test(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.PathValue("path"))
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}

type Entry struct {
	Name  string `json:"name"`
	IsDir bool   `json:"is_dir"`
}

func listContent(w http.ResponseWriter, r *http.Request) {
	subPath := r.PathValue("path")
	entries, err := os.ReadDir(path.Join(rootDirectory, subPath))
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	files := make([]Entry, 0)
	for _, entry := range entries {
		files = append(files, Entry{
			Name:  entry.Name(),
			IsDir: entry.IsDir(),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(files)
	if err != nil {
		log.Println(err)
		return
	}
}
