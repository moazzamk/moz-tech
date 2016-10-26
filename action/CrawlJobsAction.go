package action

import (
	"github.com/moazzamk/moz-tech/structures"
	"github.com/moazzamk/moz-tech/crawler"
	"github.com/moazzamk/moz-tech/service"
)

type CrawlJobsAction struct {
	salaryParser *service.SalaryParser
	skillParser *service.SkillParser
	dateParser *service.DateParser
	config  *structures.Dictionary
	storage service.Storage

	jobWriter chan structures.JobDetail
}

func NewCrawlJobsAction(
	salaryParser *service.SalaryParser,
	skillParser *service.SkillParser,
	dateParser *service.DateParser,
	config *structures.Dictionary,
	storage service.Storage) *CrawlJobsAction {

	jobWriter := make(chan structures.JobDetail)

	ret := new(CrawlJobsAction)
	ret.salaryParser = salaryParser
	ret.skillParser = skillParser
	ret.dateParser = dateParser
	ret.jobWriter = jobWriter
	ret.storage = storage
	ret.config = config

	return ret
}

func (r *CrawlJobsAction) Run() {
	go func(chan structures.JobDetail) {
		for i := range r.jobWriter {
			r.storage.AddJob(i)
			//fmt.Println(i)
		}
	}(r.jobWriter)

	var doneChannels []chan bool
	//doneChannels = append(doneChannels, startRemoteWork(&client, jobDetailWriter))
	//doneChannels = append(doneChannels, startLinkedIn(&client, jobDetailWriter))
	doneChannels = append(doneChannels, r.startStackOverflow())
	//doneChannels = append(doneChannels, r.startDice())


	for i := range doneChannels {
		_ = <-doneChannels[i]
	}
}

func (r *CrawlJobsAction) startDice() chan bool {

	var doneChannel = make(chan bool)

	go func(doneChannel chan bool) {
		worker := crawler.NewDiceCrawler(r.salaryParser, r.skillParser, r.dateParser)
		worker.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json?pgcnt=500&text=python`
		worker.JobWriter = r.jobWriter
		worker.Storage = r.storage

		worker.Crawl()

		close(doneChannel)
	}(doneChannel)

	return doneChannel
}

func (r *CrawlJobsAction) startStackOverflow() chan bool {
	var doneChannel = make(chan bool)

	go func(doneChannel chan bool, jobWriter chan structures.JobDetail) {
		worker := new(crawler.StackOverflow)
		worker.Url = `http://stackoverflow.com/jobs`
		worker.Host = `http://stackoverflow.com/`
		worker.JobWriter = jobWriter
		worker.Storage = r.storage

		worker.Crawl()

		close(doneChannel)
	}(doneChannel, r.jobWriter)

	return doneChannel
}


/*
func (r *CrawlJobsAction) startLinkedIn() chan bool {
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
*/

/*
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
*/
