package main

import (
	"io"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

const HEADER_CONTENT_TYPE = "Content-Type"

func TestCrawler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawler Suite")
}

type Client struct {
	URL        string
	HTTPClient *http.Client
}

func NewClient(URL string) *Client {
	return &Client{URL, &http.Client{}}
}

func (c *Client) DoAPIGet(url string) (*http.Response, error) {
	return c.DoAPICall("GET", c.URL+url, nil)
}

func (c *Client) DoAPIPost(url string, body io.Reader) (*http.Response, error) {
	return c.DoAPICall("POST", c.URL+url, body)
}

func (c *Client) DoAPICall(method, url string, body io.Reader) (*http.Response, error) {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set(HEADER_CONTENT_TYPE, "application/json")
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

var _ = BeforeSuite(func() {
	Setup()
})

var _ = AfterSuite(func() {
})
