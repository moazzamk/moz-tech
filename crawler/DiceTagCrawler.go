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

type DiceTagCrawler struct {
	skillWriter chan string
	Skills []string
	Url    string

	skillParser *service.SkillParser
}

func NewDiceTagCrawler(skillWriter chan string, skillParser *service.SkillParser) *DiceTagCrawler {
	dice := new(DiceTagCrawler)
	dice.skillParser = skillParser
	dice.skillWriter = skillWriter

	return dice
}

// Crawl() starts the crawling process. It is the only method anyone outside this object cares about
func (dice *DiceTagCrawler) Crawl() {
	jobChannel := make(chan string)
	url := dice.Url

	// Start routines for getting job details
	for i := 0; i < 1; i++ {
		go dice.getDetails(jobChannel)
	}

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
			jobChannel <- detailUrl
		}

		if rs[`nextUrl`] == nil {
			break
		}

		nextUrl = rs[`nextUrl`].(string)
		rs = dice.fetchSearchResults(`http://service.dice.com` + nextUrl)
		fmt.Println(`search results came back`)
	}
}

func (dice *DiceTagCrawler) fetchSearchResults(url string) map[string]interface{} {
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

func (dice *DiceTagCrawler) getDetails(jobChannel chan string) {
	for url := range jobChannel {
		doc, err := goquery.NewDocument(url)
		if err != nil {
			log.Fatal(err)
			fmt.Println(err, "ERRRR")
		}

		_ = dice.getJobSkill(doc)
	}
}

func (dice *DiceTagCrawler) getJobSkill(doc *goquery.Document) []string {
	var sss string

	doc.Find(`#labelskill`).Each(func(i int, s *goquery.Selection) {
		sss = s.Text()
	})

	uniqueSlice := structures.NewUniqueSlice(strings.Split(sss, `,`))
	dice.skillParser.ParseFromTags(uniqueSlice)

	return uniqueSlice.ToSlice()
}

