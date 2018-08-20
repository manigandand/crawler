package main

import (
	"io"
	"net/http"
	"net/url"
	"strings"

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

func (c *Client) DoAPIPost(endpoint string, website string) (*http.Response, error) {
	form := url.Values{}
	form.Add("website_address", website)

	req, _ := http.NewRequest(http.MethodPost, c.URL+endpoint, strings.NewReader(form.Encode()))
	req.Header.Set(HEADER_CONTENT_TYPE, "application/x-www-form-urlencoded")

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
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
