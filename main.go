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

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", HandleIndex)
	crawler := rtr.PathPrefix("/crawler").Subrouter()
	crawler.HandleFunc("/", spiderman).Methods(http.MethodPost)
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

func respondError(w http.ResponseWriter, err error, statusCode int) {
	http.Error(w, err.Error(), statusCode)
	return
}
