package main

import (

	"github.com/moazzamk/moz-tech/app"
	"github.com/moazzamk/moz-tech/service"
	"gopkg.in/olivere/elastic.v3"
	"github.com/moazzamk/moz-tech/action"
	"github.com/moazzamk/moz-tech/structures"
	"fmt"
	"flag"
)

/*

- find jobs that pay the most
- find companies that pay the most
- find technologies that relate to each other

*/

func main() {
	fmt.Println(`Cli started`)

	command := flag.String(`cmd`, ``, `hi there`);
	config := moz_tech.NewAppConfig(`config/config.txt`)
	esUrl, _ := config.Get(`es_url`)
	flag.Parse()

	client, err := elastic.NewClient(
		elastic.SetURL(esUrl),
		elastic.SetMaxRetries(10),
		elastic.SetSniff(false),
		elastic.SetScheme(`https`))

	if err != nil {
		panic(err)
	}

	fmt.Println(`Elastic client initialized`)

	storage := service.NewStorage(client)
	skillParser := service.NewSkillParser(storage)
	salaryParser := service.SalaryParser{}
	dateParser := service.DateParser{}


	if *command == `del-index` {
		fmt.Println(`Starting delete index`)
		(action.NewTruncateIndexAction(client)).Run(`jobs`)

	} else if *command == `crawl-jobs` {
		fmt.Println(`Starting jobs`)
		crawlJobs(&salaryParser, &skillParser, &dateParser, config, storage)

	} else if *command == `crawl-tags` {
		crawlTags(&skillParser, config, storage)

	}
}

func crawlJobs(
	salaryParser *service.SalaryParser,
	skillParser *service.SkillParser,
	dateParser *service.DateParser,
	config *structures.Dictionary,
	storage service.Storage,
) {
	action := action.NewCrawlJobsAction(
		salaryParser,
		skillParser,
		dateParser,
		config,
		storage)
	action.Run()
}


func crawlTags(skillParser *service.SkillParser, config *structures.Dictionary, storage service.Storage) {
	crawlAction := action.NewCrawlTagsAction(skillParser, config, storage)
	crawlAction.Run()
}
