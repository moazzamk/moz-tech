package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/structures"
)

type Dice struct {
	JobWriter chan structures.JobDetail
	Storage *service.Storage
	Skills []string
	Url    string

	salaryParser *service.SalaryParser
	skillParser *service.SkillParser
	dateParser *service.DateParser
}

func NewDiceCrawler(
	salaryParser *service.SalaryParser,
	skillParser *service.SkillParser,
	dateParser *service.DateParser) *Dice {

	dice := new(Dice)
	dice.salaryParser = salaryParser
	dice.skillParser = skillParser
	dice.dateParser = dateParser

	return dice
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (dice *Dice) Crawl() {
	jobChannel := make(chan structures.JobDetail)
	url := dice.Url

	// Start routines for getting job details
	for i := 0; i < 5; i++ {
		go dice.getDetails(jobChannel)
	}

	fmt.Println(`URL: ` + url)

	rs := dice.fetchSearchResults(url)
	fmt.Println(`search results came back with `, rs["count"].(float64), " results")

	if rs[`lastDocument`].(float64) <= 0 {
		fmt.Println(`No jobs found`)
		return
	}

	detailUrl := ``
	nextUrl := ``
	for rs[`resultItemList`] != nil {
		items := rs[`resultItemList`].([]interface{})
		for _, item := range items {
			obj := item.(map[string]interface{})
			detailUrl = obj[`detailUrl`].(string)

			tmp := structures.JobDetail{}
			tmp.Link = detailUrl
			if dice.Storage.HasJobWithUrl(detailUrl) {
				fmt.Println(`DICE JOB INDEXED ALREADY ` + detailUrl)
				continue
			}

			jobChannel <- tmp
		}

		if rs[`nextUrl`] == nil {
			break
		}

		nextUrl = rs[`nextUrl`].(string)
		rs = dice.fetchSearchResults(`http://service.dice.com` + nextUrl)
		fmt.Println(`search results came back`)
	}
}

func (dice *Dice) getJobDescription(doc *goquery.Document) string {
	var ret string

	doc.Find(`#jobdescSec`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (dice *Dice) fetchSearchResults(url string) map[string]interface{} {
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

func (dice *Dice) getDetails(jobChannel chan structures.JobDetail) {
	for job := range jobChannel {
		doc, err := goquery.NewDocument(job.Link)
		if err != nil {
			log.Fatal(err)
			fmt.Println(err, "ERRRR")
		}

		var ret structures.JobDetail
		ret.Telecommute, ret.Travel = dice.getTelecommuteAndTravel(doc)
		ret.Description = dice.getJobDescription(doc)
		ret.PostedDate = dice.getPostedDate(doc)
		ret.Salary = dice.getSalaryRange(doc)
		ret.Employer = dice.getEmployer(doc)
		ret.Location = dice.getLocation(doc)
		ret.Skills = dice.getJobSkill(doc)
		ret.JobType = dice.getJobType(doc)
		ret.Title = dice.getJobTitle(doc)
		ret.Source = `dice.com`

		dice.JobWriter <- ret
	}
}

func (dice *Dice) getLocation(doc *goquery.Document) string {
	var ret string
	doc.Find(`.location`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

func (dice *Dice) getEmployer(doc *goquery.Document) string {
	var ret string
	doc.Find(`.employer .dice-btn-link`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

/*
Get salary from the job posting and translate it to yearly salary
if the salary isnt already yearly
*/
func (dice *Dice) getSalaryRange(doc *goquery.Document) (*structures.SalaryRange) {
	var salaryParser service.SalaryParser
	var salary string

	doc.Find(`.icon-bank-note`).Each(func (i int, s *goquery.Selection) {
		s.Parent().Siblings().Find(`.mL20`).Each(func(j int, r *goquery.Selection) {
			salary = r.Text()
		})
	})

	return salaryParser.Parse(salary)
}

func (dice *Dice) getJobTitle(doc *goquery.Document) string {
	var title string

	doc.Find(`#jt`).Each(func (i int, s *goquery.Selection) {
		title = s.Text()
	})

	return title
}

func (dice *Dice) getTelecommuteAndTravel(doc *goquery.Document) (int, int) {
	telecommute := 0
	travel := 0

	doc.Find(`.icon-network-2`).Each(func (i int, s *goquery.Selection) {
		content := s.Parent().Siblings().Text()

		if !strings.Contains(content, `Telecommuting not available`) {
			telecommute = 1
		}

		if strings.Contains(content, `Travel`) && !strings.Contains(content, `Travel not`) {
			travel = 1
		}
	})

	return telecommute, travel
}

func (dice *Dice) getPostedDate(doc *goquery.Document) string {
	var ret string

	doc.Find(`.posted`).Each(func (i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return dice.dateParser.Parse(ret)
}

func (dice *Dice) getJobType(doc *goquery.Document) []string {
	var ret []string

	doc.Find(`.icon-briefcase`).Each(func (i int, s *goquery.Selection) {
		s.Parent().Siblings().Find(`.iconsiblings span`).Each(func(j int, r *goquery.Selection) {
			ret = strings.Split(r.Text(), `,`)
		})
	})

	return ret
}

func (dice *Dice) getJobSkill(doc *goquery.Document) []string {
	var sss string

	doc.Find(`#labelskill`).Each(func(i int, s *goquery.Selection) {
		sss = s.Text()
	})

	uniqueSlice := structures.NewUniqueSlice(strings.Split(sss, `,`))
	dice.skillParser.ParseFromTags(uniqueSlice)

	description := dice.getJobDescription(doc)
	uniqueSlice = uniqueSlice.Merge(dice.skillParser.ParseFromDescription(description))

	return uniqueSlice.ToSlice()
}

func (dice *Dice) stopWord(subject string) {

}
