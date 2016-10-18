package crawler

import (
	"github.com/moazzamk/moz-tech/structures"
	"gopkg.in/olivere/elastic.v3"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"encoding/json"
)

type LinkedIn struct {
	JobWriter chan structures.JobDetail
	Search **elastic.Client
	Skills []string
	Url    string
}

func (l *LinkedIn) Crawl() {
	doc, _ := goquery.NewDocument(l.Url)

	var response map[string]interface{}
	var jsonString string

	jsonChannel := make(chan map[string]interface{})

	// 10 threads for processing details
	for i := 0; i < 10; i++ {
		go l.getDetails(jsonChannel, l.JobWriter)
	}

	doc.Find(`#decoratedJobPostingsModule`).Each(func (i int, s *goquery.Selection) {
		jsonString, _ = s.Html()
		start := strings.Index(jsonString, `<!--`)
		end := strings.Index(jsonString, `-->`)
		jsonString = jsonString[start+4:end]
	})

	err := json.Unmarshal([]byte(jsonString), &response)
	if err != nil {
		panic(err)
	}

	items := response[`elements`].([]interface{})
	for _, item := range items {
		jsonChannel <- item.(map[string]interface{})
	}
}

func (l *LinkedIn) getDetails(jsonChannel chan map[string]interface{}, jobWriter chan structures.JobDetail) {
	var response map[string]interface{}
	var jsonString string

	for _ = range jsonChannel {
		obj := <- jsonChannel

		url := obj["viewJobCanonicalUrl"].(string)
		doc, _ := goquery.NewDocument(url)

		doc.Find(`#jobDescriptionModule`).Each(func(i int, s *goquery.Selection) {
			jsonString, _ = s.Html()
			start := strings.Index(jsonString, `<!--`)
			end := strings.Index(jsonString, `-->`)
			jsonString = jsonString[start+4:end]
		})

		err := json.Unmarshal([]byte(jsonString), &response)
		if err != nil {
			panic(err)
		}

		decoratedListing := obj["decoratedJobPosting"].(map[string]interface{})
		jobDetail := structures.JobDetail{}
		jobDetail.Source = `linkedin.com`
		jobDetail.Link = url
		jobDetail.Title = decoratedListing["jobPosting"].(map[string]interface{})["title"].(string)
		jobDetail.Employer = decoratedListing["companyName"].(string)
		jobDetail.Description = l.getDescription(doc)
		//jobDetail.JobType = ""
		jobDetail.Location = ""
		jobDetail.PostedDate = ""
		//jobDetail.Salary = ""
		//jobDetail.Skills = ""

		jobWriter <- jobDetail
	}
}

func (l *LinkedIn) getDescription(doc *goquery.Document) string {
	var ret string
	doc.Find(`meta[property="og:description"]`).Each(func(i int, s *goquery.Selection) {
		ret, _ = s.Attr(`content`)
	})

	return ret
}
