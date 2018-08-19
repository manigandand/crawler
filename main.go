package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
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

	content, err := scrapPage("https://www.redhat.com/")
	if err != nil {
		log.Fatal(err)
	}

	input := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		log.Println("could not scrap the page. err: ", err.Error())
		return
	}

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		// fmt.Println(item.Text())
		href, ok := item.Attr("href")
		if ok {
			fmt.Println(href)
		}
	})

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

// scrapPage takes url scraps the page return html string and error
func scrapPage(urlStr string) (string, error) {
	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
