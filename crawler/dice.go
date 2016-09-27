package crawler

import (
	"fmt"
	"net/http"
	"regexp"
	"io/ioutil"
)

type Dice struct {
	Url string
}


func (i *Dice) getSalaryRange(content string ) []string {
	return new ([2]string)
}

func (dice *Dice) Crawl() {

	url := dice.Url

	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error")
	}

	defer resp.Body.Close();

	//re := regexp.MustCompile("https://www.dice.com/jobs/detail/[^\\\"]+")
	re := regexp.MustCompile("https://www.dice.com/jobs/[^\\\"]+")

	body, _ := ioutil.ReadAll(resp.Body)
	matches := re.FindAllString(string(body), 100)


	for i := 0; i < len(matches); i++ {
		dice.getDetails(matches[i])
	}

	//fmt.Println("response" , string(body), matches)
}


func (i *Dice) getDetails(url string) {
	req, err := http.NewRequest("GET", url, nil)

	client := &http.Client{}
	resp, err := client.Do(req)
	if (err != nil) {
		fmt.Println("Error parsing", url)
	}

	fmt.Println(resp.Body)
}

func (i *Dice) getJobType(content string) string {
	
	return ""
}

