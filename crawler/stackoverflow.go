package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/structures"
	"gopkg.in/olivere/elastic.v3"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var wg sync.WaitGroup
var skillMutex sync.Mutex
var largestSalary float64
var largestLink string

type StackOverflow struct {
	JobWriterChannel chan structures.JobDetail
	Search           **elastic.Client
	Skills           []string
	Url              string
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (so *StackOverflow) Crawl() {
	url := so.Url
	fmt.Println(`URL: ` + url)

	ret := make(map[string]int)
	rs := so.fetchSearchResults(url)
	fmt.Println(`search results came back with `, rs["count"].(float64), " results")

	if rs[`lastDocument`].(float64) <= 0 {
		fmt.Println(`No jobs found`)
		return
	}

	detailUrl := ``
	nextUrl := ``
	for rs[`resultItemList`] != nil {
		items := rs[`resultItemList`].([]interface{})
		wg.Add(len(items))
		for _, item := range items {
			obj := item.(map[string]interface{})
			detailUrl = obj[`detailUrl`].(string)

			// Start a go routine to get details of the page
			go func(myUrl string) {
				fmt.Println(`details start for` + myUrl)
				jobDetails := so.getDetails(myUrl)
				fmt.Println(`details came back for` + myUrl)
				for i := 0; i < len(jobDetails.Skills); i++ {
					tmp := strings.ToLower(jobDetails.Skills[i])

					skillMutex.Lock()
					if _, ok := ret[tmp]; ok {
						ret[tmp]++
					} else {
						ret[tmp] = 1
					}
					skillMutex.Unlock()
				}
				wg.Done()
			}(detailUrl)
		}

		wg.Wait()

		fmt.Println("Higest salary", largestSalary, " ", largestLink)

		sortedKeys := SortedKeys(ret)
		for _, k := range sortedKeys {
			fmt.Println(k, ret[k])
		}

		if rs[`nextUrl`] == nil {
			break
		}

		nextUrl = rs[`nextUrl`].(string)
		rs = so.fetchSearchResults(`http://service.so.com` + nextUrl)
		fmt.Println(`search results came back`)
	}
}

func (so *StackOverflow) getJobDescription(doc *goquery.Document) string {
	var ret string

	doc.Find(`#jobdescSec`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
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

func (so *StackOverflow) getDetails(url string) structures.JobDetail {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	salaryRange := so.getSalaryRange(doc)
	if salaryRange.MaxSalary > largestSalary {
		largestSalary = salaryRange.MaxSalary
		largestLink = url
	}
	if salaryRange.CalculatedMaxYearlySalary > largestSalary {
		largestSalary = salaryRange.CalculatedMaxYearlySalary
		largestLink = url
	}
	if salaryRange.Salary > largestSalary {
		largestSalary = salaryRange.Salary
		largestLink = url
	}

	var ret structures.JobDetail
	ret.Telecommute, ret.Travel = so.getTelecommuteAndTravel(doc)
	ret.Description = so.getJobDescription(doc)
	ret.PostedDate = so.getPostedDate(doc)
	ret.Employer = so.getEmployer(doc)
	ret.Location = so.getLocation(doc)
	ret.Skills = so.getJobSkill(doc)
	ret.JobType = so.getJobType(doc)
	ret.Salary = salaryRange

	fmt.Println(ret)

	return ret
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
	doc.Find(`.employer .so-btn-link`).Each(func(i int, s *goquery.Selection) {
		ret = s.Text()
	})

	return ret
}

/*
Get salary from the job posting and translate it to yearly salary
if the salary isnt already yearly
*/
func (so *StackOverflow) getSalaryRange(doc *goquery.Document) *structures.SalaryRange {
	ret := new(structures.SalaryRange)
	doc.Find(`.icon-bank-note`).Each(func(i int, s *goquery.Selection) {

		str := s.Parent().Siblings().Text()
		re := regexp.MustCompile(`[$0-9,.kK]+\s*(-|to)*\s*[$0-9,.kK]+`)
		charsToReplace := map[string]string{
			`k`:  `000`,
			`K`:  `000`,
			`,`:  ``,
			`$`:  ``,
			`to`: `-`,
			` `:  ``,
		}

		ret.OriginalSalary = str
		tmp := re.FindString(str)

		if tmp == `` {
			fmt.Println(str, " was empty")
			return
		}

		for j, v := range charsToReplace {
			tmp = strings.Replace(tmp, j, v, -1)
		}

		rangeArray := strings.Split(tmp, `-`)
		rangeArrayLen := len(rangeArray)
		if rangeArrayLen == 2 {
			ret.MinSalary, _ = strconv.ParseFloat(rangeArray[0], 64)
			ret.MaxSalary, _ = strconv.ParseFloat(rangeArray[1], 64)
		} else if rangeArrayLen == 1 { // Salary is not a range
			ret.Salary, _ = strconv.ParseFloat(rangeArray[0], 64)
		}

		// Calculate yearly salary if its an hourly position
		if strings.Contains(str, `hr`) || strings.Contains(str, `hour`) {
			ret.CalculatedMinYearlySalary = ret.MinSalary * 40 * 52
			ret.CalculatedMaxYearlySalary = ret.MaxSalary * 40 * 52
			ret.CalculatedSalary = ret.Salary * 40 * 52
		}

		return
	})

	return ret
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
	travel := 0

	doc.Find(`.icon-network-2`).Each(func(i int, s *goquery.Selection) {
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

func (so *StackOverflow) getJobType(doc *goquery.Document) string {
	var ret string

	doc.Find(`.icon-briefcase`).Each(func(i int, s *goquery.Selection) {
		ret = s.Parent().Siblings().Text()
	})

	return ret
}

func (so *StackOverflow) getJobSkill(doc *goquery.Document) []string {
	var sss string

	doc.Find(`#labelskill`).Each(func(i int, s *goquery.Selection) {
		sss = s.Text()
	})

	uniqueSlice := structures.NewUniqueSlice(strings.Split(sss, `,`))

	fmt.Println("SL", uniqueSlice, sss)

	skills := so.processJobSkill(uniqueSlice)

	// Extract skills from description
	description := so.getJobDescription(doc)
	description = strings.ToLower(description)
	descriptionSentences := strings.Split(description, `. `)
	for i := range descriptionSentences {
		tmpSkill := make(map[string]int)
		tmp := strings.Split(descriptionSentences[i], ` `)
		for j := range tmp {
			tmp1 := strings.Trim(strings.Replace(tmp[j], `,`, ` `, -1), ` `)
			if len([]rune(tmp1)) >= 3 {
				tmp1 = strings.Trim(so.getNormalizedSkillSynonym(tmp1), ` `)
				tmpSkill[tmp1] = 1
			}
		}

		for j := range tmpSkill {
			if !strings.Contains(j, ` `) && service.SearchHasSkill(so.Search, j) {
				skills.Append(j)
			}
		}
	}

	return skills.ToSlice()
}

func (so *StackOverflow) stopWord(subject string) {

}

func (so *StackOverflow) processJobSkill(skills *structures.UniqueSlice) *structures.UniqueSlice {
	ret := skills

	for index, value := range skills.ToSlice() {
		tmp := strings.ToLower(strings.Trim(value, ` `))
		tmp = so.getNormalizedSkillSynonym(tmp)
		ret.Set(index, tmp)

		// If skill is more than 1 word, then check if it has multiple skills listed
		tmpSlice := strings.Split(tmp, ` `)
		tmpSliceLen := len(tmpSlice)
		for i := range tmpSlice {
			searchHasSkill := service.SearchHasSkill(so.Search, tmpSlice[i])
			if searchHasSkill {
				ret.Append(tmpSlice[i])
			}
		}

		// If the skill is one word and not present in our storage then add it

		searchHasSkill := service.SearchHasSkill(so.Search, tmp)

		if tmpSliceLen == 1 && !searchHasSkill {
			_, err := service.SearchAddSkill(so.Search, tmp)
			if err != nil {
				panic(err)
			}
			fmt.Println(`Added skill ` + tmp)
		}
	}

	return ret
}

// Correct all spellings, etc of the skill and normalize synonyms to 1 name
func (so *StackOverflow) getNormalizedSkillSynonym(skill string) string {
	ret := skill
	synonyms := map[string][]string{
		`mongo`: []string{
			`mongodb`,
			`mongo db`,
		},
		`redhat`: []string{
			`red hat`,
		},
		`javascript`: []string{
			`java script`,
			`jafascript`,
		},
		`angular`: []string{
			`angularjs`,
			`angular.js`,
			`angular js`,
		},
		`ember`: []string{
			`ember.js`,
			`emberjs`,
		},
		`mysql`: []string{
			`my sql`,
		},
		`mssql`: []string{
			`sql server`,
			`ms server`,
		},
		`aws`: []string{
			`amazon web services`,
		},
		`java`: []string{
			`corejava`,
			`core java`,
			`java8`,
		},
		`nodejs`: []string{
			`node js`,
			`node.js`,
		},
		`bootstrap`: []string{
			`boot strap`,
		},
		`bigdata`: []string{
			`big data`,
		},
		`elasticsearch`: []string{
			`elastic search`,
		},
		`machine_learning`: []string{
			`machine learning`,
		},
		`cognitive_computing`: []string{
			`cognitive computing`,
		},
		`cloud_computing`: []string{
			`cloud computing`,
		},
		`data_warehouse`: []string{
			`data warehouse design`,
			`data warehouse`,
			`data warehousing`,
		},
		`automated_testing`: []string{
			`automation test`,
		},
		`data_mining`: []string{
			`data mining`,
		},

		`predictive_analytics`: []string{
			`predictive analytics`,
		},
		`version_control`: []string{
			`version control`,
			`vcs`,
		},
		`business_intelligence`: []string{
			`business_intelligence`,
			` bi `,
			`bi `,
		},
		`azure`: []string{
			`ms azure`,
		},
		`business_analysis`: []string{
			`business analysis`,
			`business analyst`,
		},
		`data_science`: []string{
			`data science`,
			`data scientist`,
		},
	}
	for key, values := range synonyms {
		for i := range values {
			ret = strings.Replace(ret, values[i], key, -1)
		}
	}

	return ret
}
