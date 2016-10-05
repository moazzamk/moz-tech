package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"gopkg.in/olivere/elastic.v3"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sort"
	"github.com/moazzamk/moz-tech/service"
	"sync"
	"strconv"
	"time"
	"github.com/moazzamk/moz-tech/structures"
)

var wg sync.WaitGroup
var mutex sync.Mutex
var skillMutex sync.Mutex
var largestSalary float64
var largestLink string

type Dice struct {
	Url    string
	Skills []string
	Search **elastic.Client
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (dice *Dice) Crawl() {
	url := dice.Url
	fmt.Println(`URL: ` + url)

	ret := make(map[string]int)
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
		wg.Add(len(items))
		for _, item := range items {
			obj := item.(map[string]interface{})
			detailUrl = obj[`detailUrl`].(string)

			// Start a go routine to get details of the page
			go func (myUrl string) {
				fmt.Println(`details start for` + myUrl)
				jobDetails := dice.getDetails(myUrl)
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

func (dice *Dice) processJobSkill(skills []string) []string {
	ret := skills

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
		`aws`: []string {
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

	for index := range skills {

		tmp := strings.ToLower(strings.Trim(skills[index], ` `))

		// Correct all spellings, etc of the skill and normalize synonyms to 1 name
		for key, values := range synonyms {
			for i := range values {
				tmp = strings.Replace(tmp, values[i], key, -1)
			}
		}

		ret[index] = tmp

		// If skill is more than 1 word, then check if it has multiple skills listed
		tmpSlice := strings.Split(tmp, ` `)
		tmpSliceLen := len(tmpSlice)
		for i := 0; i < tmpSliceLen && tmpSliceLen > 0; i++ {
			mutex.Lock()
			searchHasSkill := service.SearchHasSkill(dice.Search, tmpSlice[i])
			mutex.Unlock()
			if  searchHasSkill {
				ret = append(ret, tmpSlice[i])
			}
		}

		// If the skill is one word and not present in our storage then add it

		mutex.Lock()
		searchHasSkill := service.SearchHasSkill(dice.Search, tmp)
		mutex.Unlock()

		if tmpSliceLen == 1 && ! searchHasSkill {
			_, err := service.SearchAddSkill(dice.Search, tmp)
			if err != nil {
				panic(err)
			}
			fmt.Println(`Added skill ` + tmp)
		}
	}


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

func (dice *Dice) getDetails(url string) structures.JobDetail {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	salaryRange := dice.getSalaryRange(doc)
	if (salaryRange.MaxSalary > largestSalary) {
		largestSalary = salaryRange.MaxSalary
		largestLink = url
	}
	if (salaryRange.CalculatedMaxYearlySalary > largestSalary) {
		largestSalary = salaryRange.CalculatedMaxYearlySalary
		largestLink = url
	}
	if (salaryRange.Salary > largestSalary) {
		largestSalary = salaryRange.Salary
		largestLink = url
	}

	var ret structures.JobDetail
	ret.Telecommute, ret.Travel = dice.getTelecommuteAndTravel(doc)
	ret.Description = dice.getJobDescription(doc)
	ret.PostedDate = dice.getPostedDate(doc)
	ret.Employer = dice.getEmployer(doc)
	ret.Location = dice.getLocation(doc)
	ret.Skills = dice.getJobSkill(doc)
	ret.JobType = dice.getJobType(doc)
	ret.Salary = salaryRange

	fmt.Println(ret)

	return ret
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
	ret := new(structures.SalaryRange)
	doc.Find(`.icon-bank-note`).Each(func (i int, s *goquery.Selection) {

		str := s.Parent().Siblings().Text()
		re := regexp.MustCompile(`[$0-9,.kK]+\s*(-|to)*\s*[$0-9,.kK]+`)
		charsToReplace := map[string]string{
			`k`: `000`,
			`K`: `000`,
			`,`: ``,
			`$`: ``,
			`to`: `-`,
			` `: ``,
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
		if (strings.Contains(str, `hr`) || strings.Contains(str, `hour`)) {
			ret.CalculatedMinYearlySalary = ret.MinSalary * 40 * 52
			ret.CalculatedMaxYearlySalary = ret.MaxSalary * 40 * 52
			ret.CalculatedSalary = ret.Salary * 40 * 52
		}

		return
	})

	return ret
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

		if strings.Contains(ret, `ago`) {
			re := regexp.MustCompile(`[0-9]+`)
			match := re.FindString(ret)
			sub, err := strconv.Atoi(match)
			if err != nil {
				ret = `Error parsing date ` + match
			}

			ts := time.Now()
			if strings.Contains(ret, `day`) {
				ts = ts.AddDate(0, 0, -1 * sub)

			} else if strings.Contains(ret, `week`) {
				ts = ts.AddDate(0, 0, -7 * sub)

			} else {
				ts = ts.AddDate(0, -1 * sub, 0)
			}

			ret = ts.String()
		}
	})

	return ret
}

func (dice *Dice) getJobType(doc *goquery.Document) string {
	var ret string

	doc.Find(`.icon-briefcase`).Each(func (i int, s *goquery.Selection) {
		ret = s.Parent().Siblings().Text()
	})

	return ret
}

func (dice *Dice) getJobSkill(doc *goquery.Document) []string {
	var sss string

	doc.Find(`#labelskill`).Each(func(i int, s *goquery.Selection) {
		sss = s.Text()
	})

	skills := dice.processJobSkill(strings.Split(sss, `,`))

	// Extract skills from description
	description := dice.getJobDescription(doc)
	descriptionSentences := strings.Split(description, `. `)
	for i := range descriptionSentences {
		tmp := strings.Split(descriptionSentences[i], ` `)
		for j := range tmp {
			tmp1 := strings.Trim(strings.Replace(tmp[j], `,`, ``, -1), ` `)
			fmt.Println(`Looking up skill` + tmp1)
			if service.SearchHasSkill(dice.Search, tmp1) {
				skills = append(skills, tmp[j])
			}
			fmt.Println(`Done Looking up skill` + tmp1)
		}
	}

	return skills
}

type SortedMap struct {
	m map[string]int
	s []string
}

func (sm *SortedMap) Len() int {
	return len(sm.m)
}

func (sm *SortedMap) Less(i, j int) bool {
	return sm.m[sm.s[i]] > sm.m[sm.s[j]]
}

func (sm *SortedMap) Swap(i, j int) {
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

func SortedKeys(m map[string]int) []string {
	sm := new(SortedMap)
	sm.m = m
	sm.s = make([]string, len(m))

	i := 0
	for key, _ := range m {
		sm.s[i] = key
		i++
	}

	sort.Sort(sm)
	return sm.s
}
