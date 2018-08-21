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
	Name      string `json:"name"`
	URL       string `json:"url"`
	ChildURLs []*URL `json:"child_urls"`
}

var (
	// SiteMapLock mutex
	SiteMapLock sync.RWMutex
	// SiteMap holds all the links in the given site
	SiteMap = make(map[string][]*URL)
)

func spiderman(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	websiteAddress := strings.TrimSpace(r.PostFormValue("website_address"))
	fmt.Println("Crawling... ", websiteAddress)
	baseURL, fullpath, err := getBaseURL(websiteAddress)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}

	// check if the site already scraped
	if siteMap, ok := readSiteMap(fullpath); ok {
		response := traverseAllLinks(removeCircularDataStructures(siteMap))
		respondSuccess(w, map[string]interface{}{"data": removeCircularDataStructures(response)})
		return
	}

	urls, err := scraper(websiteAddress, baseURL)
	if err != nil {
		respondError(w, err, http.StatusInternalServerError)
		return
	}
	urls = removeCircularDataStructures(urls)
	writeSiteMap(fullpath, urls)
	response := removeCircularDataStructures(traverseAllLinks(urls))

	respondSuccess(w, map[string]interface{}{"data": response})
	return
}

func traverseAllLinks(urls []*URL) (response []*URL) {
	var wg sync.WaitGroup

	maxNbConcurrentGoroutines := 50
	concurrentGoroutines := make(chan struct{}, maxNbConcurrentGoroutines)
	for i := 0; i < maxNbConcurrentGoroutines; i++ {
		concurrentGoroutines <- struct{}{}
	}
	done := make(chan bool)
	totalURLs := len(urls)
	resultChan := make(chan *URL)
	fmt.Println("Total Links found: ", totalURLs)

	go func() {
		for i := 0; i < totalURLs; i++ {
			<-done
			// fmt.Printf("%d from %d links\n", i+1, totalURLs)
			concurrentGoroutines <- struct{}{}
		}
	}()

	wg.Add(totalURLs)
	for _, u := range urls {
		<-concurrentGoroutines
		// fmt.Println(u.URL)
		go webCrawler(u, &wg, done, resultChan)
	}

	for i := 0; i < totalURLs; i++ {
		res := <-resultChan
		// fmt.Printf("got %d response from %d links\n", i+1, totalURLs)
		if res != nil && res.ChildURLs != nil {
			response = append(response, res)
		}
	}

	wg.Wait()
	close(done)
	close(resultChan)

	return
}

func webCrawler(u *URL, wg *sync.WaitGroup, done chan bool, resultChan chan *URL) {
	defer func() {
		wg.Done()
		done <- true
		resultChan <- u
	}()

	if childURLs, ok := readSiteMap(u.URL); ok {
		u.ChildURLs = childURLs
		return
	}

	baseURL, fullpath, err := getBaseURL(u.URL)
	if err != nil {
		return
	}

	urls, err := scraper(u.URL, baseURL)
	if err != nil {
		return
	}

	// append childURLs
	u.ChildURLs = urls
	writeSiteMap(fullpath, urls)
	return
}

func readSiteMap(fullpath string) ([]*URL, bool) {
	SiteMapLock.RLock()
	siteMap, ok := SiteMap[fullpath]
	SiteMapLock.RUnlock()

	return siteMap, ok
}

func writeSiteMap(fullpath string, urls []*URL) {
	// save in map
	SiteMapLock.Lock()
	SiteMap[fullpath] = urls
	SiteMapLock.Unlock()
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

func removeCircularDataStructures(data []*URL) (res []*URL) {
	dupLinks := map[string]bool{}
	for _, u := range data {
		if _, ok := dupLinks[u.URL]; ok {
			continue
		}
		dupLinks[u.URL] = true
		res = append(res, u)
	}

	return
}

func respondSuccess(w http.ResponseWriter, data interface{}) {
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
	if url.Scheme == "" && url.Host == "" && len(url.Path) > 2 {
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

	defer response.Body.Close()
	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	return string(data), nil
}
