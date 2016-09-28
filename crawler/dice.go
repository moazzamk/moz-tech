package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"io/ioutil"
	"strings"
)

type Dice struct {
	Url string
}

/*
func (i *Dice) getSalaryRange(content string ) []string {
	return new ([2]string)
}
*/

func (dice *Dice) Crawl() {
	dat, _ := ioutil.ReadFile(`./hourly_pay.html`)
	content := string(dat)

	fmt.Println("skill", dice.getJobSkill(content))
;
	text, tel, yo  := dice.getTelecommute(content);
	fmt.Println(text, tel, yo);

	fmt.Println(dice.getPostedDate(content))



/*
	url := dice.Url

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println("Search request complete")

	defer resp.Body.Close();

	//re := regexp.MustCompile("https://www.dice.com/jobs/detail/[^\\\"]+")
	re := regexp.MustCompile("https://www.dice.com/jobs/d[^\\ ]+")

	body, _ := ioutil.ReadAll(resp.Body)
	content := string(body)
	content = strings.Replace(content, "\r", ` `, -1)
	content = strings.Replace(content, "\n", ` `, -1)

	matches := re.FindAllString(content, 1000)

	for i := 0; i < len(matches); i++ {
		dice.getDetails(matches[i])
	}
*/
}


func (i *Dice) getDetails(url string) {
	req, err := http.NewRequest("GET", url, nil)


	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error")
	}

	defer resp.Body.Close();

	body, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("Details for " + url)

	content := string(body)
	content = strings.Replace(content, "\r", ` `, -1)
	content = strings.Replace(content, "\n", ` `, -1)

	i.getJobSkill(content)
}

func (dice *Dice) getTelecommute(content string) (string, int, int) {
	re := regexp.MustCompile(`class="mL20">[^<]+`)
	text := re.FindString(content)
	telLower := strings.ToLower(text)

	telecommute := 1
	if strings.Contains(telLower, `telecommuting not avail`) {
		telecommute = 0
	}

	travel := 0
	if strings.Contains(telLower, `travel`) && !strings.Contains(telLower, `travel not`) {
		travel = 1
	}

	return text, telecommute, travel
}

func (dice *Dice) getPostedDate(content string) string {
	re := regexp.MustCompile(`Posted Date: </b>[^<]+`)
	date := re.FindString(content)
	date = strings.Replace(date, `Posted Date: </b>`, ``, -1)
	date = strings.Trim(date, ` `)

	return date
}

func (i *Dice) getJobSkill(content string) []string {
	re := regexp.MustCompile(`labelskill" title = "[^"]+`)

	match := re.FindString(content)
	match = strings.Replace(match, `labelskill" title = "`, ``, -1)

	skills := strings.Split(match, `,`)

	return skills
}

func (i *Dice) getJobType(content string) string {


	return ""
}

