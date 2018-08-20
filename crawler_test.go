package main

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Crawler", func() {
	Context("Test scrapPage", func() {
		BeforeEach(func() {
			SiteMap = nil
		})
		AfterEach(func() {
			SiteMap = nil
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
			url := "https://www.redhat.com/"
			htmlStr, err := scrapPage(url)
			Ω(err).ShouldNot(HaveOccurred())
			Expect(htmlStr).NotTo(Equal(""))
		})
	})

	Context("Test getBaseURL", func() {
		BeforeEach(func() {
			SiteMap = nil
		})
		AfterEach(func() {
			SiteMap = nil
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
			SiteMap = nil
		})
		AfterEach(func() {
			SiteMap = nil
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

	It("Should return full path", func() {
		// var w http.ResponseWriter
		// respondSuccess(w, nil)
		// url := "https://www.redhat.com/about/?q=123445"
		// baseURL := "https://www.redhat.com"
		// fullpath, err := validateURL(url, baseURL)
		// Ω(err).ShouldNot(HaveOccurred())
		// Expect(fullpath).To(Equal("https://www.redhat.com/about/"))
	})

})
