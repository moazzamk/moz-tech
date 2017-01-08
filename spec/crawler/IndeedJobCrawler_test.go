package crawler

import (
	. "github.com/moazzamk/moz-tech/mock/service"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"github.com/moazzamk/moz-tech/crawler"
	"github.com/golang/mock/gomock"

	"net/http/httptest"
	"net/http"
	"io/ioutil"
	"fmt"
)

var mockCtrl *gomock.Controller

func TestIndeedJobCrawler(t *testing.T) {
	mockCtrl = gomock.NewController(t)

	RegisterFailHandler(Fail)
	RunSpecs(t, "Indeed crawler Suite")
}

var _ = Describe("Indeed job crawler", func () {
	Context("Search listings", func () {
		It("crawls multiple pages of job listings", func () {
			var lastPage string

			server := httptest.NewServer(http.HandlerFunc(func(rs http.ResponseWriter, rq *http.Request) {
				lastPage = rq.URL.Query().Get("start")
				if rq.RequestURI == `/?` {
					content, err := ioutil.ReadFile(`../../fixtures/indeed_list.html`)
					if err != nil {
						fmt.Println("FL", err)
					}

					rs.Write(content)
				}
			}))
			defer server.Close()

			salaryParser := NewMockISalaryParser(mockCtrl)
			skillParser := NewMockISkillParser(mockCtrl)
			dateParser := NewMockIDateParser(mockCtrl)

			crawler := crawler.NewIndeedJobCrawler(
				salaryParser,
				skillParser,
				dateParser)

			crawler.Url = server.URL + "/?"

			crawler.Crawl()

			Expect(lastPage).To(Equal(`1`))
		})
	})
})
