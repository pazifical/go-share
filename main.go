package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	_ "embed"
)

type Directory struct {
	Name    string
	Entries []Entry
}

type Entry struct {
	Name      string
	Directory string
	IsDir     bool
}

var port = 8910

//go:embed static
var static embed.FS

var templ *template.Template

func main() {
	server, err := createServer()
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), server))
}

func createServer() (*http.ServeMux, error) {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("."))

	staticFS, err := fs.Sub(static, "static")
	if err != nil {
		return nil, err
	}
	staticFileserver := http.FileServer(http.FS(staticFS))

	mux.Handle("/", http.StripPrefix("/", staticFileserver))
	mux.Handle("/api/download/", http.StripPrefix("/api/download/", fileServer))
	mux.HandleFunc("/api/list/{path...}", renderEntries)

	return mux, nil
}

func init() {
	t, err := template.ParseFS(static, "static/data.html")
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

func GetContentType(fPath string) (string, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	buf := make([]byte, 512)
	_, err = f.Read(buf)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buf), nil

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

	err = templ.Execute(w, Directory{Name: subPath, Entries: entries})
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
