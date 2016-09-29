package main 

import (
	"./crawler"
)

/*

- find jobs that pay the most
- find companies that pay the most
- find technologies that relate to each other

 */

func main() {
	diceCrawler := new(crawler.Dice)
	diceCrawler.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json`
	diceCrawler.Crawl()



	/*
	linkedInCrawler := new (crawler.LinkedIn)
	linkedInCrawler.Crawl()
	*/
}
