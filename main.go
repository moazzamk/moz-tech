package main 

import (
	`./crawler`
	`gopkg.in/olivere/elastic.v3`
	"github.com/moazzamk/moz-tech/app"
)

/*

- find jobs that pay the most
- find companies that pay the most
- find technologies that relate to each other

 */



func main() {
	config := moz_tech.NewAppConfig(`config/config.txt`)
	esUrl, _ := config.Get(`es_url`)

	// Initialize ES
	// Keep sniffing to false. It causes ES library to fail
	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetMaxRetries(10),
		elastic.SetSniff(false),
		elastic.SetScheme(`https`))
	if (err != nil) {
		panic(err)
	}

	// Make sure ES works
	_, _ , err = client.Ping(esUrl).Do()
	if err != nil {
		panic(err)
	}

	diceCrawler := new(crawler.Dice)
	diceCrawler.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json?pgcnt=500&text=phalcon&state=TX`
	diceCrawler.Search = &client
	diceCrawler.Crawl()


	/*
	linkedInCrawler := new (crawler.LinkedIn)
	linkedInCrawler.Crawl()
	*/
}
