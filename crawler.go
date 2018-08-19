package main

import (
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

// URL holds the basic site links
type URL struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	// IsVisited bool   `json:"is_visited"`
}

var (
	SiteMapLock sync.RWMutex
	// SiteMap holds all the links in the given site
	SiteMap = make(map[string][]*URL)
)

func spiderman(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	websiteAddress := strings.TrimSpace(r.PostFormValue("website_address"))
	baseURL, fullpath, err := getBaseURL(websiteAddress)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	// check if the site already scraped
	SiteMapLock.Lock()
	if siteMap, ok := SiteMap[fullpath]; ok {
		SiteMapLock.Unlock()
		respondSuccess(w, siteMap)
		return
	}
	SiteMapLock.Unlock()

	urls, err := scraper(websiteAddress, baseURL)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	// save in map
	SiteMapLock.Lock()
	SiteMap[fullpath] = urls
	SiteMapLock.Unlock()
	go traverseAllLinks(urls)

	respondSuccess(w, urls)
	return
}

func traverseAllLinks(urls []*URL) {
	var wg sync.WaitGroup

	maxNbConcurrentGoroutines := 5
	concurrentGoroutines := make(chan struct{}, maxNbConcurrentGoroutines)
	for i := 0; i < maxNbConcurrentGoroutines; i++ {
		concurrentGoroutines <- struct{}{}
	}
	done := make(chan bool)
	totalURLs := len(urls)

	go func() {
		for i := 0; i < totalURLs; i++ {
			<-done
			concurrentGoroutines <- struct{}{}
		}
	}()

	wg.Add(totalURLs)
	for _, u := range urls {
		<-concurrentGoroutines
		go webCrawler(u, &wg)
	}

	wg.Wait()
	close(done)
	return
}

func webCrawler(u *URL, wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		// SiteMapLock.Unlock()
	}()

	SiteMapLock.Lock()
	if _, ok := SiteMap[u.URL]; !ok {
		SiteMapLock.Unlock()
		baseURL, fullpath, err := getBaseURL(u.URL)
		if err != nil {
			return
		}
		urls, err := scraper(u.URL, baseURL)
		if err != nil {
			return
		}

		// save in map
		SiteMapLock.Lock()
		SiteMap[fullpath] = urls
		SiteMapLock.Unlock()
	}

	return
}

func scraper(websiteAddress, baseURL string) ([]*URL, error) {
	content, err := scrapPage(websiteAddress)
	if err != nil {
		return nil, err
	}

	input := strings.NewReader(string(content))
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		log.Println("could not scrap the page. err: ", err.Error())
		return nil, err
	}
	var urls []*URL

	doc.Find("a").Each(func(index int, item *goquery.Selection) {
		// fmt.Println(item.Text())
		href, ok := item.Attr("href")
		if ok {
			subURL, err := validateURL(strings.TrimSpace(href), baseURL)
			if err == nil {
				urls = append(urls, &URL{
					Name: strings.TrimSpace(item.Text()),
					URL:  strings.TrimSpace(subURL),
				})
			}
			// fmt.Println(href)
		}
	})

	return urls, nil
}

func respondSuccess(w http.ResponseWriter, data []*URL) {
	gz := gzip.NewWriter(w)
	defer gz.Close()

	buf, err := json.Marshal(data)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)
	if _, err := gz.Write(buf); err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	return
}

func validateURL(urlStr, baseURL string) (string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	if url.Scheme == "" && url.Host == "" {
		return fmt.Sprintf("%s%s", baseURL, url.Path), nil
	}

	if strings.Contains(urlStr, baseURL) {
		return fmt.Sprintf("%s%s", baseURL, url.Path), nil
	}

	return "", errors.New("invalid url")
}

func getBaseURL(urlStr string) (string, string, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return "", "", err
	}

	return fmt.Sprintf("%s://%s", url.Scheme, url.Host),
		fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path), nil
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
