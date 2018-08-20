package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crawler", func() {
	Context("Test scrapPage", func() {
		BeforeEach(func() {
			SiteMap = map[string][]*URL{}
		})
		AfterEach(func() {
			SiteMap = map[string][]*URL{}
		})
		It("Should Fail - Invalid url to parse", func() {
			url := "http://test.com/Segment%%2815197306101420000%29.ts"
			htmlStr, err := scrapPage(url)
			Expect(err.Error()).To(Equal(fmt.Sprintf(`parse %s: invalid URL escape "%s"`, url, "%%2")))
			Expect(htmlStr).To(Equal(""))
		})
		It("Should Fail - Invalid url to parse", func() {
			url := "http://[fe80::1%25en0]:/"
			htmlStr, err := scrapPage(url)
			Expect(err.Error()).To(Equal("Get http://[fe80::1%25en0]/: dial tcp [fe80::1%en0]:80: connect: invalid argument"))
			Expect(htmlStr).To(Equal(""))
		})
		It("Should Fail - Invalid url to parse", func() {
			url := "https://www.google.com/"
			htmlStr, err := scrapPage(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(htmlStr).NotTo(Equal(""))
		})
	})

	Context("Test getBaseURL", func() {
		BeforeEach(func() {
			SiteMap = map[string][]*URL{}
		})
		AfterEach(func() {
			SiteMap = map[string][]*URL{}
		})
		It("Should Fail - Invalid url to parse", func() {
			url := "http://test.com/Segment%%2815197306101420000%29.ts"
			baseURL, fullpath, err := getBaseURL(url)
			Expect(err.Error()).To(Equal(fmt.Sprintf(`parse %s: invalid URL escape "%s"`, url, "%%2")))
			Expect(baseURL).To(Equal(""))
			Expect(fullpath).To(Equal(""))
		})

		It("Should return base url and full path", func() {
			url := "https://www.redhat.com/about/?q=123445"
			baseURL, fullpath, err := getBaseURL(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(baseURL).To(Equal("https://www.redhat.com"))
			Expect(fullpath).To(Equal("https://www.redhat.com/about/"))
		})
	})

	Context("Test validateURL", func() {
		BeforeEach(func() {
			SiteMap = map[string][]*URL{}
		})
		AfterEach(func() {
			SiteMap = map[string][]*URL{}
		})
		It("Should Fail - Invalid url to parse", func() {
			url := "http://test.com/Segment%%2815197306101420000%29.ts"
			baseURL := "https://www.redhat.com"
			fullpath, err := validateURL(url, baseURL)
			Expect(err.Error()).To(Equal(fmt.Sprintf(`parse %s: invalid URL escape "%s"`, url, "%%2")))
			Expect(fullpath).To(Equal(""))
		})

		It("Should append url scheme and host and return full path", func() {
			url := "/about/?q=123445"
			baseURL := "https://www.redhat.com"
			fullpath, err := validateURL(url, baseURL)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fullpath).To(Equal("https://www.redhat.com/about/"))
		})

		It("Should Fail - url not belongs to the same domain", func() {
			url := "https://www.google.com/about/?q=123445"
			baseURL := "https://www.redhat.com"
			fullpath, err := validateURL(url, baseURL)
			Expect(err.Error()).To(Equal("invalid url"))
			Expect(fullpath).To(Equal(""))
		})

		It("Should return full path", func() {
			url := "https://www.redhat.com/about/?q=123445"
			baseURL := "https://www.redhat.com"
			fullpath, err := validateURL(url, baseURL)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(fullpath).To(Equal("https://www.redhat.com/about/"))
		})
	})

	Context("POST http://127.0.0.1:8080/crawler/", func() {
		It("POST crawler - invalid url", func() {
			url := "http://test.com/Segment%%2815197306101420000%29.ts"

			res, err := tClient.Crawler(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusInternalServerError))
		})
		It("POST crawler - should return urls", func() {
			url := "https://www.redhat.com/"

			res, err := tClient.Crawler(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			var data []*URL
			jsonErr := json.NewDecoder(res.Body).Decode(&data)
			Ω(jsonErr).ShouldNot(HaveOccurred())
			// fmt.Println(data)
			Expect(len(data)).NotTo(Equal(0))
		})
		It("POST crawler - should fetch from the cache", func() {
			url := "https://www.redhat.com/"

			res, err := tClient.Crawler(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(res.StatusCode).To(Equal(http.StatusOK))
			var data []*URL
			jsonErr := json.NewDecoder(res.Body).Decode(&data)
			Ω(jsonErr).ShouldNot(HaveOccurred())
			// fmt.Println(data)
			Expect(len(data)).NotTo(Equal(0))
		})
	})

})

func (c *Client) Crawler(urlStr string) (*http.Response, error) {
	return c.DoAPIPost("/crawler/", urlStr)
}
