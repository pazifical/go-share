package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	_ "embed"
)

type Entry struct {
	Name      string
	Directory string
	IsDir     bool
}

var port = 8910

//go:embed static/index.html
var index []byte

//go:embed static/data.html
var data embed.FS

var templ *template.Template

func init() {
	t, err := template.ParseFS(data, "static/data.html")
	if err != nil {
		log.Fatal(err)
	}
	templ = t
}

func getDirEntries(directory string) ([]Entry, error) {
	dirEntries, err := os.ReadDir(directory)
	if err != nil {
		return []Entry{}, err
	}

	files := make([]Entry, 0)
	for _, entry := range dirEntries {
		files = append(files, Entry{
			Name:      entry.Name(),
			Directory: directory,
			IsDir:     entry.IsDir(),
		})
	}
	return files, nil
}

func main() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/api/download/", http.StripPrefix("/api/download/", fs))

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/api/list/{path...}", renderEntries)

	fmt.Printf("Open a browser on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func renderEntries(w http.ResponseWriter, r *http.Request) {
	subPath := r.PathValue("path")
	if subPath == "" {
		subPath = "."
	}

	entries, err := getDirEntries(subPath)
	if err != nil {
		log.Println(err)
		entries = make([]Entry, 0)
	}

	err = templ.Execute(w, entries)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Write(index)
}
