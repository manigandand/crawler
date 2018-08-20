package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	tServer *httptest.Server
	tClient *Client
)

func TestMain(m *testing.M) {
	MainSetup()
	defer MainTearDown()
	os.Exit(m.Run())
}

func MainSetup() {
	// SiteMap = make(map[string][]*URL)
	InitAPI()

	if tServer == nil {
		tServer = httptest.NewServer(rtr)
	}
}

func MainTearDown() {
	SiteMap = make(map[string][]*URL)
	if tServer != nil {
		tServer.Close()
	}
}

func Setup() {
	if tClient == nil {
		tClient = NewClient(tServer.URL)
	}
}

var _ = Describe("API Test", func() {
	Describe("GET http://127.0.0.1:8080/", func() {
		Context("GET http://127.0.0.1:8080/", func() {
			It("GET crawler index", func() {
				res, err := tClient.CrawlerHomePage()
				Ω(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(http.StatusOK))
			})
		})
		Context("GET http://127.0.0.1:8080/crawler/status/", func() {
			It("GET crawler status endpoint - expect zero result", func() {
				res, err := tClient.CrawlerStatus()
				Ω(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(http.StatusOK))
				var data []*URL
				jsonErr := json.NewDecoder(res.Body).Decode(&data)
				Ω(jsonErr).ShouldNot(HaveOccurred())
				Expect(len(data)).To(Equal(0))
			})

			It("GET crawler status endpoint - expect two results", func() {
				urls := []*URL{
					&URL{
						Name: "Red Hat About",
						URL:  "https://www.redhat.com/about/",
					},
					&URL{
						Name: "Red Hat Contact",
						URL:  "https://www.redhat.com/contact/",
					},
					&URL{
						Name: "Red Hat Products",
						URL:  "https://www.redhat.com/en/topics/cloud-computing/why-choose-red-hat-cloud",
					},
				}

				SiteMapLock.Lock()
				SiteMap = map[string][]*URL{
					"https://www.redhat.com/":         urls,
					"https://www.redhat.com/about/":   urls,
					"https://www.redhat.com/contact/": urls,
				}
				SiteMapLock.Unlock()

				res, err := tClient.CrawlerStatus()
				Ω(err).ShouldNot(HaveOccurred())
				Expect(res.StatusCode).To(Equal(http.StatusOK))

				var fr map[string][]*URL
				jsonErr := json.NewDecoder(res.Body).Decode(&fr)
				Ω(jsonErr).ShouldNot(HaveOccurred())
				Expect(len(fr)).To(Equal(3))
			})
		})
	})

})

func (c *Client) CrawlerHomePage() (*http.Response, error) {
	return c.DoAPIGet("/")
}

func (c *Client) CrawlerStatus() (*http.Response, error) {
	return c.DoAPIGet("/crawler/status/")
}
