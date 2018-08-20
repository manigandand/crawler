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
			url := "http://test.com/Segment%%2815197306101420000%29.ts"
			htmlStr, err := scrapPage(url)
			Expect(err.Error()).To(Equal(fmt.Sprintf(`parse %s: invalid URL escape "%s"`, url, "%%2")))
			Expect(htmlStr).To(Equal(""))
		})
	})
})
