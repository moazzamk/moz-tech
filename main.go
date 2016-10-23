package main

import (
	"./crawler"
	"fmt"
	"github.com/moazzamk/moz-tech/app"
	"github.com/moazzamk/moz-tech/service"
	"github.com/moazzamk/moz-tech/structures"
	"gopkg.in/olivere/elastic.v3"
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

	if err != nil {
		panic(err)
	}

	fmt.Println(client, &client)

	// Make sure ES works
	_, _, err = client.Ping(esUrl).Do()
	if err != nil {
		panic(err)
	}

	storage := service.NewStorage(&client)
	skillParser := service.NewSkillParser(&storage)
	salaryParser := service.SalaryParser{}
	dateParser := service.DateParser{}


	fmt.Println("Initialized")

	jobDetailWriter := make(chan structures.JobDetail)

	go func(chan structures.JobDetail) {
		for i := range jobDetailWriter {
			fmt.Println(i)
			storage.AddJob(i)
		}
	}(jobDetailWriter)

	var doneChannels []chan bool
	//doneChannels = append(doneChannels, startRemoteWork(&client, jobDetailWriter))
	//doneChannels = append(doneChannels, startLinkedIn(&client, jobDetailWriter))
	doneChannels = append(doneChannels, startDice(&storage, &salaryParser, &skillParser, &dateParser, jobDetailWriter))
	doneChannels = append(doneChannels, startStackOverflow(&storage, jobDetailWriter))

	for i := range doneChannels {
		_ = <-doneChannels[i]
	}
}

func startDice(
	storage *service.Storage,
	salaryParser *service.SalaryParser,
	skillParser *service.SkillParser,
	dateParser *service.DateParser,
	jobDetailWriter chan structures.JobDetail) chan bool {

	var doneChannel = make(chan bool)

	go func(doneChannel chan bool) {
		worker := crawler.NewDiceCrawler(salaryParser, skillParser, dateParser)
		worker.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json?pgcnt=500&text=python`
		worker.JobWriter = jobDetailWriter
		worker.Storage = storage

		worker.Crawl()

		close(doneChannel)
	}(doneChannel)

	return doneChannel
}

func startLinkedIn(client **elastic.Client, jobWriter chan structures.JobDetail) chan bool {
	var doneChannel = make(chan bool)

	go func(doneChannel chan bool, jobWriter chan structures.JobDetail) {
		worker := new(crawler.LinkedIn)
		worker.Url = `https://www.linkedin.com/jobs/search?keywords=&location=Dallas%2FFort%20Worth%20Area&locationId=`
		worker.JobWriter = jobWriter
		worker.Search = client
		worker.Crawl()

		//time.Sleep(1000 * time.Millisecond)


		close(doneChannel)
	}(doneChannel, jobWriter)

	return doneChannel
}


func startStackOverflow(storage *service.Storage, JobDetailWriter chan structures.JobDetail) chan bool {
	var doneChannel = make(chan bool)

	go func (doneChannel chan bool, jobDetailWriter chan structures.JobDetail) {
		worker := new(crawler.StackOverflow)
		worker.Url = `http://stackoverflow.com/jobs`
		worker.Host = `http://stackoverflow.com/`
		worker.JobWriter = JobDetailWriter
		worker.Storage = storage

		worker.Crawl()

		close (doneChannel)
	}(doneChannel, JobDetailWriter)

	return doneChannel
}

func startRemoteWork(client **elastic.Client, JobDetailWriter chan structures.JobDetail) chan bool {
	var doneChannel = make(chan bool)

	fmt.Println(`Started remotework.io crawler`)
	go func(doneChannel chan bool, searchWriteChannnel chan structures.JobDetail) {
		remoteWorkOkCrawler := new(crawler.RemoteWork)
		remoteWorkOkCrawler.Url = `https://remoteok.io/index.json`
		remoteWorkOkCrawler.Host = `https://remoteok.io`
		remoteWorkOkCrawler.Search = client
		remoteWorkOkCrawler.SearchWriteChannel = JobDetailWriter
		remoteWorkOkCrawler.Crawl()
		doneChannel <- true

		close(doneChannel)
	}(doneChannel, JobDetailWriter)

	return doneChannel
}
