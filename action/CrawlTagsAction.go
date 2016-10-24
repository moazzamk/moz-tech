package action

import (
	"github.com/moazzamk/moz-tech/structures"
	"github.com/moazzamk/moz-tech/crawler"
	"github.com/moazzamk/moz-tech/service"
)

type CrawlTagsAction struct {
	skillParser *service.SkillParser
	config  *structures.Dictionary
	storage *service.Storage

	skillWriter chan string
}

func NewCrawlTagsAction(
	skillParser *service.SkillParser,
	config *structures.Dictionary,
	storage *service.Storage) *CrawlTagsAction {

	skillWriter := make(chan string)

	ret := new(CrawlTagsAction)
	ret.skillParser = skillParser
	ret.skillWriter = skillWriter
	ret.storage = storage
	ret.config = config

	return ret
}

func (r *CrawlTagsAction) Run() {
	go func(chan string) {
		for i := range r.skillWriter {
			r.storage.AddSkill(i)
		}
	}(r.skillWriter)

	var doneChannels []chan bool
	//doneChannels = append(doneChannels, startRemoteWork(&client, jobDetailWriter))
	//doneChannels = append(doneChannels, startLinkedIn(&client, jobDetailWriter))
	doneChannels = append(doneChannels, r.startStackOverflow())
	doneChannels = append(doneChannels, r.startDice())

	for i := range doneChannels {
		_ = <-doneChannels[i]
	}
}

func (r *CrawlTagsAction) startDice() chan bool {

	var doneChannel = make(chan bool)

	go func(doneChannel chan bool) {
		worker := crawler.NewDiceTagCrawler(r.skillWriter, r.skillParser)
		worker.Url = `http://service.dice.com/api/rest/jobsearch/v1/simple.json?pgcnt=500&text=python`
		worker.Crawl()

		close(doneChannel)
	}(doneChannel)

	return doneChannel
}

func (r *CrawlTagsAction) startStackOverflow() chan bool {
	var doneChannel = make(chan bool)

	go func(doneChannel chan bool, skillWriter chan string) {
		worker := crawler.NewStackOverflowTagCrawler(r.skillWriter, r.skillParser)
		worker.Url = `http://stackoverflow.com/jobs`
		worker.Host = `http://stackoverflow.com/`
		worker.Crawl()

		close(doneChannel)
	}(doneChannel, r.skillWriter)

	return doneChannel
}


/*
func (r *CrawlTagsAction) startLinkedIn() chan bool {
	var doneChannel = make(chan bool)

	go func(doneChannel chan bool, skillWriter chan structures.JobDetail) {
		worker := new(crawler.LinkedIn)
		worker.Url = `https://www.linkedin.com/jobs/search?keywords=&location=Dallas%2FFort%20Worth%20Area&locationId=`
		worker.SkillWriter = skillWriter
		worker.Search = client
		worker.Crawl()

		//time.Sleep(1000 * time.Millisecond)

		close(doneChannel)
	}(doneChannel, skillWriter)

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
