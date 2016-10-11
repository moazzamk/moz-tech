package crawler

import (
	"gopkg.in/olivere/elastic.v3"
	"github.com/PuerkitoBio/goquery"
	"log"
	"fmt"
	"github.com/moazzamk/moz-tech/structures"
	"strings"
	"github.com/moazzamk/moz-tech/service"
)

type RemoteWork struct {
	Host string
	Url string
	Search **elastic.Client
}

func (rw *RemoteWork) Crawl() {
	var links []string

	doc, err := goquery.NewDocument(rw.Url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	doc.Find(`a[itemprop=url]`).Each(func (i int, s *goquery.Selection) {
		a, _ := s.Attr(`href`)
		links = append(links, a)
		fmt.Println(a)
	})
}

func (rw *RemoteWork) getDetails(uri string) {
	url := rw.Host + uri
	doc, err := goquery.NewDocument(rw.Host + url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	var jobDetails structures.JobDetail
	jobDetails.Link = url
	jobDetails.Description = rw.getJobDescription(doc)
	jobDetails.Employer = rw.getEmployer(doc)
	jobDetails.JobType = rw.getJobType(doc)
	jobDetails.Location = rw.getLocation(doc)
	jobDetails.PostedDate = rw.getPostedDate(doc)
	jobDetails.Salary = rw.getSalaryRange(doc)
	jobDetails.Title = rw.getJobTitle(doc)
	jobDetails.Skills = rw.getSkills(doc)
	jobDetails.Telecommute = rw.getTelecommuteAndTravel(doc)
	jobDetails.Travel = rw.getTelecommuteAndTravel(doc)

	service.SearchAddJob(rw.Search, jobDetails)
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

	doc.Find(``).Each(func (i int, s *goquery.Selection) {
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

func (rw *RemoteWork) getPostedDate(doc *goquery.Document) {
	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
}

func (rw *RemoteWork) getJobType(doc *goquery.Document) {
	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
}

func (rw *RemoteWork) getTelecommuteAndTravel(doc *goquery.Document) {
	doc.Find(``).Each(func (i int, s *goquery.Selection) {

	})
}
