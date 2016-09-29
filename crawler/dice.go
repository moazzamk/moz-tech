package crawler

import (
	"fmt"
	"regexp"
	"strings"
	"../goquery"
	"log"
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
	//dat, _ := ioutil.ReadFile(`./hourly_pay.html`)
	//content := string(dat)


	dice.getDetails(`https://www.dice.com/jobs/detail/Senior-PHP-Developer-RealPage-Inc.-Richardson-TX-75082/10111244/641877?icid=sr1-1p&q=php&l=Dallas,%20TX`)

return;
/*	fmt.Println("skill", dice.getJobSkill(content))

	tel, yo  := dice.getTelecommuteAndTravel(content);
	fmt.Println(tel, ":", yo);

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


func (dice *Dice) getDetails(url string) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		fmt.Println(err, "ERRRR")
	}

	fmt.Println(dice.getJobSkill(doc))
}

func (dice *Dice) getTelecommuteAndTravel(content string) (int, int) {
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

	// @todo : process jobs posted hours ago not weeks ago

	re := regexp.MustCompile(`Posted Date: </b>[^<]+`)
	date := re.FindString(content)
	date = strings.Replace(date, `Posted Date: </b>`, ``, -1)
	date = strings.Trim(date, ` `)

	return date
}

func (dice *Dice) getJobType(doc *goquery.Document) string {
	return ``
}

func (dice *Dice) getJobSkill(doc *goquery.Document) []string {
	var sss string;

	doc.Find(`#labelskill`).Each(func (i int, s *goquery.Selection) {
		sss = s.Text();
	})

	return strings.Split(sss, `,`)
}
