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
	"time"
)

type StackOverflow struct {
	JobWriter chan structures.JobDetail
	Storage   *service.Storage
	Skills    []string
	Url       string
	Host      string
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (so *StackOverflow) Crawl() {
	var totalJobs int
	var jobs string

	jobChannel := make(chan structures.JobDetail)
	url := so.Url

	fmt.Println(`URL: ` + url)

	// Start routines for getting job details
	for i := 0; i < 5; i++ {
		go so.getDetails(so.JobWriter, jobChannel)
	}

	doc, _ := goquery.NewDocument(url)

	// Get total number of jobs
	doc.Find(`#index-hed .description`).Each(func(i int, s *goquery.Selection) {
		jobs = strings.Replace(s.Text(), " jobs", "", -1)
		jobs = strings.Replace(jobs, ",", "", -1)
		jobs = strings.Replace(jobs, ` `, ``, -1)
		totalJobs, _ = strconv.Atoi(jobs)
	})
	if totalJobs <= 0 {
		fmt.Println(`No jobs found`)
		return
	}
	fmt.Println(`Stack Overflow came back with `, totalJobs, " results")

	jobCount := so.dispatchJobs(doc, jobChannel)
	for i := 2; jobCount > 0; i++ {
		doc, _ = goquery.NewDocument(url + `?pg=` + strconv.Itoa(i))
		jobCount = so.dispatchJobs(doc, jobChannel)
	}
}

func (so *StackOverflow) dispatchJobs(doc *goquery.Document, jobChannel chan structures.JobDetail) int {
	var jobCount = 0

	doc.Find(`h2 a.job-link`).Each(func(i int, s *goquery.Selection) {
		jobCount++
		href, _ := s.Attr(`href`)

		if so.Storage.HasJobWithUrl(href) {
			fmt.Println(`SO JOB INDEXED ALREADY ` + href)
			return
		}

		job := structures.JobDetail{}
		job.PostedDate = so.getPostedDate(doc)
		job.Link = href

		jobChannel <- job
	})

	fmt.Println(`DISPATCH FINISHED WITH CNT:` + strconv.Itoa(jobCount) + ` ` + doc.Url.String())

	return jobCount
}

func (so *StackOverflow) getJobDescription(doc *goquery.Document) string {
	var ret string

	doc.Find(`.jobdetail p, .jobdetail ul`).Each(func(i int, s *goquery.Selection) {
		ret += s.Text() + "\n\n"
	})

	return ret
}

func (so *StackOverflow) getDetails(jobWriterChannel chan structures.JobDetail, jobChannel chan structures.JobDetail) {
	for job := range jobChannel {
		fmt.Println(`Starting`, job.Link)
		doc, err := goquery.NewDocument(so.Host + job.Link)
		if err != nil {
			fmt.Println(err, "ERRRR")
			log.Fatal(err)
		}

		so.Storage.HasJobWithUrl(job.Link)

		job.Telecommute, job.Travel = so.getTelecommuteAndTravel(doc)
		job.Description = so.getJobDescription(doc)
		job.Salary = so.getSalaryRange(doc)
		job.Employer = so.getEmployer(doc)
		job.Location = so.getLocation(doc)
		job.Skills = so.getJobSkill(doc)
		job.JobType = so.getJobType(doc)
		job.Source = `stackoverflow.com`

		so.JobWriter <- job
	}
}

func (so *StackOverflow) fetchSearchResults(url string) map[string]interface{} {
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


func (so *StackOverflow) getLocation(doc *goquery.Document) string {
	var ret string

	doc.Find(`.location`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (so *StackOverflow) getEmployer(doc *goquery.Document) string {
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
func (so *StackOverflow) getSalaryRange(doc *goquery.Document) *structures.SalaryRange {
	salaryParser := service.SalaryParser{}
	var salary string
	doc.Find(`.salary`).Each(func(i int, s *goquery.Selection) {
		salary = s.Text()
	})

	return salaryParser.Parse(salary)
}

func (so *StackOverflow) getJobTitle(doc *goquery.Document) string {
	var title string

	doc.Find(`#jt`).Each(func(i int, s *goquery.Selection) {
		title = s.Text()
	})

	return title
}

func (so *StackOverflow) getTelecommuteAndTravel(doc *goquery.Document) (int, int) {
	telecommute := 0
	travel := -1

	doc.Find(`.detail-remote`).Each(func (i int, s *goquery.Selection) {
		telecommute = 1
	})

	return telecommute, travel
}

func (so *StackOverflow) getPostedDate(doc *goquery.Document) string {
	var ret string

	doc.Find(`.posted`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()

		if strings.Contains(ret, `ago`) {
			re := regexp.MustCompile(`[0-9]+`)
			match := re.FindString(ret)
			sub, err := strconv.Atoi(match)
			if err != nil {
				ret = `Error parsing date ` + match
			}

			ts := time.Now()
			if strings.Contains(ret, `day`) {
				ts = ts.AddDate(0, 0, -1*sub)

			} else if strings.Contains(ret, `week`) {
				ts = ts.AddDate(0, 0, -7*sub)

			} else {
				ts = ts.AddDate(0, -1*sub, 0)
			}

			ret = ts.String()
		}
	})

	return ret
}

func (so *StackOverflow) getJobType(doc *goquery.Document) []string {
	var ret []string

	doc.Find(`.icon-briefcase`).Each(func(i int, s *goquery.Selection) {
		ret = strings.Split(s.Parent().Siblings().Text(), `,`)
	})

	return ret
}

func (so *StackOverflow) getJobSkill(doc *goquery.Document) []string {
	var tags []string
	skillParser := service.NewSkillParser(so.Storage)

	doc.Find(`.tags a.no-tag-menu`).Each(func(i int, s *goquery.Selection) {
		tags = append(tags, s.Text())
	})

	uniqueSlice := structures.NewUniqueSlice(tags)

	//fmt.Println("SL", uniqueSlice, tags)

	skills := skillParser.ParseFromTags(uniqueSlice)

	// Extract skills from description
	description := so.getJobDescription(doc)
	skills = skills.Merge(skillParser.ParseFromDescription(description))

	return skills.ToSlice()
}


