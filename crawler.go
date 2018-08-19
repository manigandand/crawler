package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func spiderman(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	websiteAddress := strings.TrimSpace(r.PostFormValue("website_address"))

	content, err := scrapPage(websiteAddress)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	input := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		log.Println("could not scrap the page. err: ", err.Error())
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		fmt.Println(item.Text())
		href, ok := item.Attr("href")
		if ok {
			fmt.Println(href)
		}
	})

	w.Write([]byte(websiteAddress))
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
	// url.Scheme
	// fmt.Println(url.Fragment)
	// fmt.Println(url.Host)
	// fmt.Println(url.Path)
	// fmt.Println(url.RawPath)
	// fmt.Println(url.RawQuery)

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
