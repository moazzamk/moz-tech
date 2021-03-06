package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/structures"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type StackOverflow struct {
	JobWriter chan structures.JobDetail
	Storage   service.Storage
	Skills    []string
	Url       string
	Host      string

	salaryParser *service.SalaryParser
	skillParser *service.SkillParser
	dateParser *service.DateParser
}

func NewStackOverflowJobCrawler(salaryParser *service.SalaryParser, skillParser *service.SkillParser, dateParser *service.DateParser) *StackOverflow {
	ret := StackOverflow{
		salaryParser: salaryParser,
		skillParser: skillParser,
		dateParser: dateParser,
	}

	return &ret
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (r *StackOverflow) Crawl() {
	var totalJobs int

	jobChannel := make(chan structures.JobDetail)
	url := r.Url

	fmt.Println(`URL: ` + url)

	// Start routines for getting job details
	for i := 0; i < 25; i++ {
		go r.getDetails(r.JobWriter, jobChannel, i)
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

func (r *StackOverflow) getTotalJobs(doc *goquery.Document) int {
	var totalJobs int

	doc.Find(`span.js-search-title`).Each(func(i int, s *goquery.Selection) {
		regex := regexp.MustCompile(`[0-9,]+`)
		jobs := regex.FindString(s.Text())
		jobs = strings.Replace(jobs, ",", "", -1)
		totalJobs, _ = strconv.Atoi(jobs)
	})

	return totalJobs
}

func (r *StackOverflow) dispatchJobs(doc *goquery.Document, jobChannel chan structures.JobDetail) int {
	var jobCount = 0

	doc.Find(`h2 a.job-link`).Each(func(i int, s *goquery.Selection) {
		jobCount++
		href, _ := s.Attr(`href`)

		/*
		if r.Storage.HasJobWithUrl(href) {
			fmt.Println(`SO JOB INDEXED ALREADY ` + href)
			return
		}
		*/

		job := structures.JobDetail{}
		job.PostedDate = r.getPostedDate(doc)
		job.Link = r.Host + href

		jobChannel <- job
	})

	fmt.Println(`DISPATCH FINISHED WITH CNT:` + strconv.Itoa(jobCount) + ` ` + doc.Url.String())

	return jobCount
}

func (r *StackOverflow) getJobDescription(doc *goquery.Document) string {
	var ret string

	doc.Find(`.jobdetail p, .jobdetail ul`).Each(func(i int, s *goquery.Selection) {
		ret += s.Text() + "\n\n"
	})

	return ret
}

func (r *StackOverflow) getDetails(jobWriterChannel chan structures.JobDetail, jobChannel chan structures.JobDetail, i int) {
	for job := range jobChannel {
		//fmt.Println(i, ` Starting`, job.Link)
		doc, err := goquery.NewDocument(job.Link)
		if err != nil {
			fmt.Println(err, "ERRRR")
			log.Fatal(err)
		}
		//fmt.Println(i, `doc came back`)

		//r.Storage.HasJobWithUrl(job.Link)

		job.Telecommute, job.Travel = r.getTelecommuteAndTravel(doc)
		job.Description = r.getJobDescription(doc)
		job.Salary = r.getSalaryRange(doc)
		job.Employer = r.getEmployer(doc)
		job.Location = r.getLocation(doc)
		job.Skills = r.getJobSkill(doc)
		job.JobType = r.getJobType(doc)
		job.Title = r.getJobTitle(doc)
		job.Source = `stackoverflow.com`

		fmt.Println(i, job.Title, `job parse finished`)

		r.JobWriter <- job
		//fmt.Println(i, ` Finsihed`, job.Link)
	}
}

func (r *StackOverflow) fetchSearchResults(url string) map[string]interface{} {
	var response map[string]interface{}

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


func (r *StackOverflow) getLocation(doc *goquery.Document) string {
	var ret string

	doc.Find(`.location`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (r *StackOverflow) getEmployer(doc *goquery.Document) string {
	var ret string

	doc.Find(`a.employer`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

/*
Get salary from the job posting and translate it to yearly salary
if the salary isnt already yearly
*/
func (r *StackOverflow) getSalaryRange(doc *goquery.Document) *structures.SalaryRange {
	var salary string
	doc.Find(`.salary`).Each(func(i int, s *goquery.Selection) {
		salary = s.Text()
	})

	return r.salaryParser.Parse(salary)
}

func (r *StackOverflow) getJobTitle(doc *goquery.Document) string {
	var title string

	doc.Find(`.detail-jobTitle`).Each(func(i int, s *goquery.Selection) {
		title = strings.Trim(
					strings.Trim(s.Text(), ` `),
					"\n")
	})

	return title
}

func (r *StackOverflow) getTelecommuteAndTravel(doc *goquery.Document) (int, int) {
	telecommute := 0
	travel := -1

	doc.Find(`.detail-remote`).Each(func (i int, s *goquery.Selection) {
		telecommute = 1
	})

	return telecommute, travel
}

func (r *StackOverflow) getPostedDate(doc *goquery.Document) string {
	var ret string

	doc.Find(`.posted`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return r.dateParser.Parse(ret)
}

func (r *StackOverflow) getJobType(doc *goquery.Document) []string {
	var ret []string

	doc.Find(`.icon-briefcase`).Each(func(i int, s *goquery.Selection) {
		ret = strings.Split(s.Parent().Siblings().Text(), `,`)
	})

	return ret
}

func (r *StackOverflow) getJobSkill(doc *goquery.Document) []string {
	var tags []string

	doc.Find(`.tags a.no-tag-menu`).Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	uniqueSlice := structures.NewUniqueSlice(tags)
	skills := r.skillParser.ParseFromTags(uniqueSlice)

	// Extract skills from description
	description := r.getJobDescription(doc)
	skills = skills.Merge(r.skillParser.ParseFromDescription(description))

	return skills.ToSlice()
}


