package main

import (
	"embed"
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
	// http.HandleFunc("/api/data", serveData)
	http.HandleFunc("/api/list/{path...}", renderEntries)
	log.Fatal(http.ListenAndServe(":8910", nil))
}

func renderEntries(w http.ResponseWriter, r *http.Request) {
	subPath := r.PathValue("path")
	if subPath == "" {
		subPath = "."
	}

	log.Printf("subPath: '%s'", subPath)
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

// func OLDlistContent(w http.ResponseWriter, r *http.Request) {
// 	subPath := r.PathValue("path")
// 	entries, err := os.ReadDir(path.Join(directory, subPath))
// 	if err != nil {
// 		log.Println(err)
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	files := make([]Entry, 0)
// 	for _, entry := range entries {
// 		files = append(files, Entry{
// 			Name:  entry.Name(),
// 			IsDir: entry.IsDir(),
// 		})
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(files)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}
// }
