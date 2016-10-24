package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"github.com/moazzamk/moz-tech/structures"
	"fmt"
	"log"
	"github.com/moazzamk/moz-tech/service"
)

type StackOverflowTagCrawler struct {
	Host string
	Url string
	skillWriter chan string
	storage service.Storage

	skillParser *service.SkillParser
}

func NewStackOverflowTagCrawler(skillWriter chan string, skillParser *service.SkillParser) *StackOverflowTagCrawler {
	ret := new(StackOverflowTagCrawler)
	ret.skillParser = skillParser
	ret.skillWriter = skillWriter
	return ret
}

func (r *StackOverflowTagCrawler) Crawl() {
	var totalJobs int

	jobChannel := make(chan string)
	url := r.Url

	fmt.Println(`URL: ` + url)

	// Start routines for getting job details
	for i := 0; i < 5; i++ {
		go r.getDetails(jobChannel)
	}

	doc, _ := goquery.NewDocument(url)

	totalJobs = r.getTotalJobs(doc)
	if totalJobs <= 0 {
		fmt.Println(`No jobs found`)
		return
	}

	fmt.Println(`Stack Overflow came back with `, totalJobs, " results")

	jobCount := r.dispatchJobs(doc, jobChannel)
	for i := 2; jobCount > 0; i++ {
		doc, _ = goquery.NewDocument(url + `?pg=` + strconv.Itoa(i))
		jobCount = r.dispatchJobs(doc, jobChannel)
	}
}

func (r *StackOverflowTagCrawler) dispatchJobs(doc *goquery.Document, jobChannel chan string) int {
	var jobCount = 0

	doc.Find(`h2 a.job-link`).Each(func(i int, s *goquery.Selection) {
		jobCount++
		href, _ := s.Attr(`href`)
		jobChannel <- href
	})

	fmt.Println(`SO DISPATCH FINISHED WITH CNT:` + strconv.Itoa(jobCount) + ` ` + doc.Url.String())

	return jobCount
}

func (r *StackOverflowTagCrawler) getDetails(jobChannel chan string) {
	for uri := range jobChannel {
		fmt.Println(`Starting`, uri)
		doc, err := goquery.NewDocument(r.Host + uri)
		if err != nil {
			fmt.Println(err, "ERRRR")
			log.Fatal(err)
		}

		skills := r.getJobSkill(doc)
		for _, skill := range skills {
			r.skillWriter <- skill
		}
	}
}

func (r *StackOverflowTagCrawler) getTotalJobs(doc *goquery.Document) int {
	var totalJobs int

	doc.Find(`#index-hed .description`).Each(func(i int, s *goquery.Selection) {
		jobs := strings.Replace(s.Text(), " jobs", "", -1)
		jobs = strings.Replace(jobs, ",", "", -1)
		jobs = strings.Replace(jobs, ` `, ``, -1)
		totalJobs, _ = strconv.Atoi(jobs)
	})

	return totalJobs
}


func (r *StackOverflowTagCrawler) getJobSkill(doc *goquery.Document) []string {
	var tags []string

	doc.Find(`.tags a.no-tag-menu`).Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	uniqueSlice := structures.NewUniqueSlice(tags)
	skills := r.skillParser.ParseFromTags(uniqueSlice)

	return skills.ToSlice()
}
