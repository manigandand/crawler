package main

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
)

var rtr *mux.Router

func InitAPI() {
	rtr = mux.NewRouter()
	rtr.HandleFunc("/", HandleIndex).Methods(http.MethodGet)
	crawler := rtr.PathPrefix("/crawler").Subrouter()
	crawler.HandleFunc("/", spiderman).Methods(http.MethodPost)
	crawler.HandleFunc("/status/", crawlerStatus).Methods(http.MethodGet)
}

func main() {
	InitAPI()

	http.Handle("/", rtr)
	log.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", nil)

	return
}

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	tmplPath := "/index.html"
	t := template.New("index.html")
	currDir, err := os.Getwd()
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	absPath := path.Clean(currDir + tmplPath)
	t, err = t.ParseFiles(absPath)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	buf := &bytes.Buffer{}
	if err := t.Execute(buf, nil); err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	io.WriteString(w, buf.String())
	return
}

func crawlerStatus(w http.ResponseWriter, r *http.Request) {
	respondSuccess(w, SiteMap)
	return
}

func respondError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
	return
}
