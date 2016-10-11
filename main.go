package main

import (
	`./crawler`
	`gopkg.in/olivere/elastic.v3`
	"github.com/moazzamk/moz-tech/app"
	"github.com/moazzamk/moz-tech/structures"
	"github.com/moazzamk/moz-tech/service"
	"fmt"
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

	fmt.Println("Initialized")

	searchWriteChannel := make(chan structures.JobDetail)

	go func (chan structures.JobDetail) {
		for i := range searchWriteChannel {
			service.SearchAddJob(&client, i)
		}
	}(searchWriteChannel)

	var doneChannels []chan bool
	doneChannels = append(doneChannels, startRemoteWork(&client))

	for i := range doneChannels {
		_ = <- doneChannels[i]
	}

	/*
	diceCrawler := new(crawler.Dice)
	diceCrawler.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json?pgcnt=500&text=phalcon&state=TX`
	diceCrawler.Search = &client
	diceCrawler.Crawl()
	*/




	/*
	linkedInCrawler := new (crawler.LinkedIn)
	linkedInCrawler.Crawl()
	*/
}

func startRemoteWork(client **elastic.Client) chan bool {
	var doneChannel = make(chan bool)

	go func (doneChannel chan bool) {
		remoteWorkOkCrawler := new(crawler.RemoteWork)
		remoteWorkOkCrawler.Url = `https://remoteok.io/remote-dev-jobs`
		remoteWorkOkCrawler.Host = `https://remoteok.io`
		remoteWorkOkCrawler.Search = client
		remoteWorkOkCrawler.Crawl()
		doneChannel <- true

		close(doneChannel)
	}(doneChannel)

	return doneChannel
}
