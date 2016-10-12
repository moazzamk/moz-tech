package crawler

import (
	"gopkg.in/olivere/elastic.v3"
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"github.com/moazzamk/moz-tech/structures"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
)


// This is not complete. I realized they have expired jobs and not enough info.
// so it's probably not worth it to index them
type RemoteWork struct {
	Host string
	Url string
	Search **elastic.Client
	SearchWriteChannel chan structures.JobDetail
}

func (rw *RemoteWork) Crawl() {
	var detailUrl string

	rs := rw.fetchSearchResults(rw.Url)

	for _, item := range rs {
		obj := item.(map[string]interface{})
		detailUrl = obj[`url`].(string)

		// Start a go routine to get details of the page
		fmt.Println(`details start for` + detailUrl)
		rw.getDetails(rw.Host + "/" + detailUrl)
		fmt.Println(`details came back for` + detailUrl)
	}
}

func (rw *RemoteWork) fetchSearchResults(url string) []interface{} {
	var response []interface{}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println(`Could not decode response`, err)
		return nil
	}

	return response
}

func (rw *RemoteWork) getDetails(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		fmt.Println(67)
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	var jobDetails structures.JobDetail
	jobDetails.Link = url
	jobDetails.Source = `remoteok.io`
	jobDetails.Description = rw.getJobDescription(doc)
	jobDetails.Employer = rw.getEmployer(doc)
	jobDetails.JobType = rw.getJobType(doc)
	jobDetails.Location = rw.getLocation(doc)
	jobDetails.PostedDate = rw.getPostedDate(doc)
	//jobDetails.Salary = rw.getSalaryRange(doc)
	//jobDetails.Title = rw.getJobTitle(doc)
	//jobDetails.Skills = rw.getSkills(doc)
	//jobDetails.Telecommute = rw.getTelecommuteAndTravel(doc)
	//jobDetails.Travel = rw.getTelecommuteAndTravel(doc)

	rw.SearchWriteChannel <- jobDetails
}

func (rw *RemoteWork) getSkills(doc *goquery.Document) []string {
	ret := new(structures.UniqueSlice)
	doc.Find(`td.tags .tag`).Each(func (i int, s *goquery.Selection) {
		attr, _ := s.Attr(`class`)
		attr = strings.Replace(attr, `tag `, ``, -1)
		attr = strings.Replace(attr, `tag-`, ``, -1)

		ret.Append(attr)
	})

	return ret.ToSlice()
}

func (rw *RemoteWork) getJobDescription(doc *goquery.Document) string {
	var ret string

	doc.Find(`.description`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (rw *RemoteWork) getLocation(doc *goquery.Document) string {
	return `Remote`
}

func (rw *RemoteWork) getEmployer(doc *goquery.Document) string {
	var ret string

	doc.Find(`a[itemprop=hiringOrganization] h3`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (rw *RemoteWork) getSalaryRange(doc *goquery.Document) {
	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
}

func (rw *RemoteWork) getJobTitle(doc *goquery.Document) string {
	var ret string

	doc.Find(`a.preventLink h2[itemprop=title]`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (rw *RemoteWork) getPostedDate(doc *goquery.Document) string {

	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
	return "test"
}

func (rw *RemoteWork) getJobType(doc *goquery.Document) string {

	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})

	return `test`
}

func (rw *RemoteWork) getTelecommuteAndTravel(doc *goquery.Document) (int, int) {

	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
	return 1, 1
}
