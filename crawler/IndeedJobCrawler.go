package crawler

import (
	"github.com/moazzamk/moz-tech/structures"
	"github.com/moazzamk/moz-tech/service"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"regexp"
	"strconv"
	"fmt"
	"log"
	"net/http"
)

type IndeedJobCrawler struct {
	JobWriter chan structures.JobDetail
	Storage   service.Storage
	Skills    []string
	Url       string
	Host      string
	Client    http.Client

	salaryParser service.ISalaryParser
	skillParser service.ISkillParser
	dateParser service.IDateParser
}

func NewIndeedJobCrawler(
	salaryParser service.ISalaryParser,
	skillParser service.ISkillParser,
	dateParser service.IDateParser) *IndeedJobCrawler {

	ret := IndeedJobCrawler{
		salaryParser: salaryParser,
		skillParser: skillParser,
		dateParser: dateParser,
	}

	return &ret
}

func (r *IndeedJobCrawler) Crawl() {
	jobChannel := make(chan structures.JobDetail)
	url := r.Url

	fmt.Println(`INDEED`, `Starting indeed crawler`)

	doc, _ := goquery.NewDocument(url)
	fmt.Println(doc.Html())
	fmt.Println(url)
	totalCount := r.getTotalCount(doc)
	jobsPerPage := r.getJobPerPage(doc)

	// Start routines for getting job details
	for i := 0; i < jobsPerPage; i++ {
		go r.getDetails(jobChannel, i)
	}

	log.Println(`INDEED`, `Total jobs:`, totalCount, jobsPerPage)

	for i := 0; i < totalCount; i += jobsPerPage {
		fmt.Println(`INDEED`, `Getting page`, i)
		url = r.Url + `&start=` + strconv.Itoa(i)
		doc, err := goquery.NewDocument(url)
		if err != nil {
			fmt.Println(err)
		}

		doc.Find(`a.turnstileLink`).Each(func (i int, s *goquery.Selection) {
			link, err := s.Attr(`href`)
			if err == false {
				fmt.Println(`INDEED ERR`, 67, err)
				return
			}

			if !strings.HasPrefix(link, `http://`) {
				link = `http://indeed.com` + link
			}

			job := structures.JobDetail{
				Link: link,
			}

			fmt.Println(`INDEED`, `Add to queue`, job.Link)
			jobChannel <- job
		})
	}
}

func  (r *IndeedJobCrawler) getJobPerPage(doc *goquery.Document) int {
	var ret int

	doc.Find(`#searchCount`).Each(func(i int, s *goquery.Selection) {
		numRegex := regexp.MustCompile(`[0-9,]+`)
		regex := regexp.MustCompile(`to [^of]+`)
		tmp := s.Text()
		str := regex.FindString(tmp)
		str = numRegex.FindString(str)
		str = strings.Replace(str, `,`, ``, -1)

		fmt.Println(`INDEEDY1`, `|` + str + `|`, tmp)
		ret, _ = strconv.Atoi(str)
	})

	return ret;
}

func  (r *IndeedJobCrawler) getTotalCount(doc *goquery.Document) int {
	var ret int

	doc.Find(`#searchCount`).Each(func(i int, s *goquery.Selection) {
		numRegex := regexp.MustCompile(`[0-9,]+`)
		regex := regexp.MustCompile(`of [^<]+`)
		tmp := s.Text()
		str := regex.FindString(tmp)
		str = numRegex.FindString(str)
		str = strings.Replace(str, `,`, ``, -1)

		fmt.Println(`INDEEDY`, `|` + str + `|`, tmp)
		ret, _ = strconv.Atoi(str)
	})

	return ret
}

func (r *IndeedJobCrawler) getDetails(jobChannel chan structures.JobDetail, index int) {
	for job := range jobChannel {
		doc, _ := goquery.NewDocument(job.Link)
		job.Title = r.getJobTitle(doc)
		if job.Title == `` {
			continue;
		}

		job.Telecommute, job.Travel = r.getTelecommuteAndTravel(doc)
		job.Description = r.getJobDescription(doc)
		job.PostedDate = r.getPostedDate(doc)
		job.Location = r.getLocation(doc)
		job.Employer = r.getEmployer(doc)
		job.JobType = r.getJobType(doc)
		job.Skills = r.getJobSkill(doc)

		job.Source = `indeed.com`

		fmt.Println(`INDEED`, index, `Writing details for `, job.Link)
		r.JobWriter <- job
		fmt.Println(`INDEED`, index, `Wrote details for `, job.Link)
	}
}

func (r *IndeedJobCrawler) getLocation(doc *goquery.Document) string {
	var ret string

	doc.Find(`.location`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (r *IndeedJobCrawler) getEmployer(doc *goquery.Document) string {
	var ret string

	doc.Find(`.company`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (r *IndeedJobCrawler) getSalaryRange(doc *goquery.Document) *structures.SalaryRange {
	var salary string
	doc.Find(`.salary`).Each(func(i int, s *goquery.Selection) {
		salary = s.Text()
	})

	return r.salaryParser.Parse(salary)
}

func (r *IndeedJobCrawler) getJobTitle(doc *goquery.Document) string {
	var title string

	doc.Find(`.jobtitle`).Each(func(i int, s *goquery.Selection) {
		title = strings.Trim(
			strings.Trim(s.Text(), ` `),
			"\n")
	})

	return title
}

func (r *IndeedJobCrawler) getTelecommuteAndTravel(doc *goquery.Document) (int, int) {
	telecommute := 0
	travel := -1

	if r.getLocation(doc) == `Remote` {
		telecommute = 1
	}

	return telecommute, travel
}

func (r *IndeedJobCrawler) getPostedDate(doc *goquery.Document) string {
	var ret string

	doc.Find(`.date`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return r.dateParser.Parse(ret)
}

func (r *IndeedJobCrawler) getJobType(doc *goquery.Document) []string {
	regex := regexp.MustCompile(`Job Type:[^<]+`)
	str := regex.FindString(doc.Text())

	return []string{strings.Replace(str, `Job Type:`, ``, -1)}
}

func (r *IndeedJobCrawler) getJobSkill(doc *goquery.Document) []string {
	description := r.getJobDescription(doc)
	skills := r.skillParser.ParseFromDescription(description)

	return skills.ToSlice()
}

func (r *IndeedJobCrawler) getJobDescription(doc *goquery.Document) string {
	var ret string
	doc.Find(`#job_summary`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}
