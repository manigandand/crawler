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

func main() {
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
