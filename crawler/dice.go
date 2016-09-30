package crawler

import (
	"fmt"
	"regexp"
	"strings"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"log"
//	"github.com/moazzamk/moz-tech/arrays"
)

type Dice struct {
	Url string
	Skills []string
}

/*
func (i *Dice) getSalaryRange(content string ) []string {
	return new ([2]string)
}
*/

func (dice *Dice) Crawl() {
	url := dice.Url + `?text=php`
	fmt.Println(`URL: ` + url)

	ret := make(map[string]int)
	rs := dice.fetchSearchResults(url)

	if (rs[`lastDocument`].(float64) <= 0) {
		fmt.Println(`No jobs found`)
		return
	}

	detailUrl := ``
	nextUrl := ``
	for rs[`resultItemList`] != nil {
		items := rs[`resultItemList`].([]interface{})
		for _,item := range items {
			obj := item.(map[string]interface{})
			detailUrl = obj[`detailUrl`].(string)

			skills := dice.getDetails(detailUrl)
			for i := 0; i < len(skills); i++ {
				tmp := strings.ToLower(skills[i])
				if _, ok := ret[tmp]; ok {
					ret[tmp]++;
				} else {
					ret[tmp] = 1
				}
			}

		}

		fmt.Println(ret)


		if rs[`nextUrl`] == nil {
			break;
		}

		fmt.Println(ret)

		nextUrl = rs[`nextUrl`].(string)
		rs = dice.fetchSearchResults(`http://service.dice.com` + nextUrl)
	}
}

func (dice *Dice) processJobSkill(skills []string) []string {
	ret := skills
/*
	syn := map[string][]string {
		`mongo`: []string{
			`mongodb`,
		},
		`redhat`: []string {
			`red hat`,

		},
		`javascript`: []string {
			`java script`,
			`js`,
		},
		`angular`: []string {
			`angularjs`,
		},
	}

*/
	skillies := []string {
		`codeigniter`,
		`laravel`,
		`zend`,
		`symfony`,
		`mvc`,
		`yii`,

		`python`,
		`ruby`,
		`mysql`,
		`postgresql`,

		`git`,

		`angular`,
		`json`,
		`javascript`,
		`jquery`,
		`rest`,


		`scrum`,
		`ajax`,
		`xml`,
		`css3`,
		`css`,

		`drupal`,
		`wordpress`,
		`joomla`,


		`agile`,
		`ruby`,
		`soap`,

		`lamp`,
	}

	// In compound skills, if we find a skill we know of, we add it to the list

	tmp := ``
	for i := 0; i < len(skillies); i++ {
		for j := 0; j < len(skills); j++ {
			tmp = strings.Trim(skills[j], ` `)
			if strings.Contains(tmp, ` `) && strings.Contains(tmp, skillies[i]) {
				ret = append(ret, skillies[i])
			}
		}
	}

	return ret
}

/*
func (dice *Dice) processJobSkill(skills []string) []string {
	var ret []string

	ignoreSkills := []string{
		`software development`,
		`development`,
	}

	for i := 0; i < len(skills); i++ {
		if exists := !arrays.InArray(ignoreSkills, skills[i]); exists {
			append(ret, skills[i])
		}
	}

	return ret
}
*/

func (dice *Dice) fetchSearchResults(url string) map[string]interface{} {
	var response map[string]interface{}

	resp, err := http.Get(url)
	if (err != nil) {
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

func (dice *Dice) getDetails(url string) []string {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	return dice.getJobSkill(doc)
}

func (dice *Dice) getTelecommuteAndTravel(content string) (int, int) {
	// This method needs to be re-implemented as we no longer pass content as string to detail methods
	telecommute := 0
	if !strings.Contains(content, `Telecommuting not available`) {
		telecommute = 1
	}

	travel := 0
	if strings.Contains(content, `Travel`) && !strings.Contains(content, `Travel not`) {
		travel = 1
	}

	return telecommute, travel
}

func (dice *Dice) getPostedDate(content string) string {
	// This method needs to be re-implemented as we no longer pass content as string to detail methods

	// @todo : process jobs posted hours ago not weeks ago

	re := regexp.MustCompile(`Posted Date: </b>[^<]+`)
	date := re.FindString(content)
	date = strings.Replace(date, `Posted Date: </b>`, ``, -1)
	date = strings.Trim(date, ` `)

	return date
}

func (dice *Dice) getJobType(doc *goquery.Document) string {
	// This method needs to be re-implemented as we no longer pass content as string to detail methods
	return ``
}

func (dice *Dice) getJobSkill(doc *goquery.Document) []string {
	var sss string;

	doc.Find(`#labelskill`).Each(func (i int, s *goquery.Selection) {
		sss = s.Text();
	})

	return dice.processJobSkill(strings.Split(sss, `,`))
}
