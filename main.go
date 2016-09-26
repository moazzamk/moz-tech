package main 

import (
	"./crawler"
)

func main() {
	diceCrawler := new(crawler.Dice)
	diceCrawler.Crawl()

	linkedInCrawler := new (crawler.LinkedIn)
	linkedInCrawler.Crawl()
}
